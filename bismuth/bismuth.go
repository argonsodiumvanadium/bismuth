package bismuth


import (
	"reflect"
	"strings"
)
/*
	_________BISMUTH_________
	
	A relational Database

	made by argonsodiumvanadium
*/

/*
	___________Logical structure___________

		First one will create a root user which
	alone has the power to create normal users,
	normal users can create users too but they 
	will be Segments as they would solely exist 
	to arrange data and will require no login 
	credentials
		Every user's database will be visible by its
	parent user or superceding user

	DATABASE structure
		Every database must have a name and here
	every database has a bunch of structures for
	efficient fetching and appending of data the 
	structures are :-
	- Tables
	  These are hashmaps in other words and are
	  purely based on hashing
	- Search Tree
	  an extremly fast tree structure that
	  takes the advantage of o(1) indexing time
	  of arrays
	- Graph
	  We all know what a graph is
*/

type (

	Session struct {
		Users map[string]*User
	}
	
	User struct {
		Name string
		Password string
		Databases map[string]*Database
		Access string
		Children []*User
		Parent *User
		Session *Session
	}

	Database struct {
		Tables map[string]*Table
		Owner *User
		SearchTreeHead *SearchTreeNode
	}
	
	SearchTreeNode struct {
		Children []*SearchTreeNode
		Values []*Entity
	}

	Table struct {
		Rows map[string]map[string][]*Entity
		Columns map[string]*Entity
	}

	Entity struct {
		Name string
		Data map[string]*Node
		Relations []*Entity
	}

	Node struct {
		Value interface{}
	}

)

const (
	ROOT = "r"
	USER = "u"
	SEGM = "s"

	NUM_OF_CHARS_SUPPORTED = 128

/*	COL_QUERY = ':'
	ROW_QUERY = '>'
	SUB_QUERY = '~'
*/
	QUERY_DELIM = ","
	MULT_ARGS = "*"
)

func CreateSession () (*Session) {
	return &Session{make(map[string]*User)}
}

func (self *Session) CreateRootUser (name,password string) (*User) {
	user := User{name,password,make(map[string]*Database),ROOT,make([]*User,0),nil,self}

	if self.Users[name] == nil {
			self.Users[name] = &(user)
	} else {
		panic("the user [ "+name+" ] has already been created")
	}

	return &user
}

//requires root access to create a user 
func (self *User) CreateUser (name,password string) (*User) {
	switch self.Access {
	case ROOT:
		user := User{name,password,make(map[string]*Database),USER,make([]*User,0),self,self.Session}
		self.Children = append(self.Children,&user)
		if self.Session.Users[name] == nil {
			self.Session.Users[name] = &(user)
		} else {
			panic("the user [ "+name+" ] has already been created")
		}
		return &user
	default:
		panic("The user is not root, only root can create users\n\tHELP : A normal user can create segments with CreateSegment() method")
	}
	return nil
}

func (self *User) CreateSegment (name string) (*User) {
	user := User{name,"",make(map[string]*Database),SEGM,make([]*User,0),self,self.Session}

	if self.Session.Users[name] == nil {
			self.Session.Users[name] = &(user)
	} else {
		panic("the user [ "+name+" ] has already been created")
	}

	self.Children = append(self.Children,&user)
	return &user
}

func (self *User) CreateDatabase (name string) {
	db := Database{ make(map[string]*Table) , self ,&SearchTreeNode{}}
	interator := self

	for interator != nil {
		interator.Databases[name] = &db
		name = interator.Name+"/"+name
		interator = interator.Parent
	}
}

func (self *Session) Login (name,passwd string) (*User) {
	user := self.Users[name]
	if user.Password == passwd {
		return user
	}
	return nil
}

func (self *User) GetDatabase (name string) (*Database) {
	return self.Databases[name]
}

func (self *User) GetAllDatabases () (map[string]*Database) {
	return self.Databases
}

func (self *Database) CreateTable (name string,fields... string) {
	self.ct(name,fields)
}

func (self *User) CreateTableIn (dbId,name string ,fields... string) {
	db := self.Databases[dbId]
	db.ct(name,fields)
}

func (self *Database) ct (name string,fields []string) {
	table := Table{ make(map[string]map[string][]*Entity) , make(map[string]*Entity) }

	for _,field := range(fields) {
		table.Rows[field] = make(map[string][]*Entity)
	}

	self.Tables[name] = &table
}

func (self *Table) AddRows (names... string) {
	for _,f := range(names) {
		self.Rows[f] = make(map[string][]*Entity)
	}
}

func (self *Database) AddEntity (entity_name string,d interface{}) {
	Fields := make([]string,0)
	name := "TABLE_WITH_FIELDS : "

	entity := Entity{entity_name,make(map[string]*Node),make([]*Entity,0)}

	data := reflect.Indirect(reflect.ValueOf(d))
	entity.Data = make(map[string]*Node)
	
	for i:=0;i<data.NumField();i++ {
		Fields = append(Fields,data.Type().Field(i).Name)
		entity.Data[Fields[i]] = &Node{data.Field(i)}
		name = name + ", " + Fields[i]
	}

	table := self.FetchTable(name)

	if table == nil {
		self.ct(name,Fields)
		table = self.FetchTable(name)
	}

	if table.Columns[name] != nil {
		panic("The given table already contains the entry "+name)
	}

	for _,val := range(Fields){
		
		if table.Rows[val][entity_name] == nil {
			table.Rows[val][entity_name] = make([]*Entity,0)
		}

		table.Rows[val][entity_name] = append(table.Rows[val][entity_name],&entity)
	}

	self.AddStringToSearch(entity_name,&entity)
}

func (self *Database) AddEntityIn (tableId ,name string,d interface{}) {
	Fields := make([]string,0)
	entity_name := name

	entity := Entity{name,make(map[string]*Node),make([]*Entity,0)}
	
	data := reflect.Indirect(reflect.ValueOf(d))
	entity.Data = make(map[string]*Node)

	for i:=0;i<data.NumField();i++{
		Fields = append(Fields,data.Type().Field(i).Name)
		entity.Data[Fields[i]] = &Node{data.Field(i)}
	}

	table := self.FetchTable(tableId)

	if table == nil {
		self.ct(name,Fields)
		table = self.FetchTable(name)
	}

	if table.Columns[name] != nil {
		panic("The given table already contains the entry "+name)
	}

	for _,val := range(Fields){
		
		if table.Rows[val] == nil {
			table.Rows[val] = make(map[string][]*Entity)
			table.Rows[val][entity_name] = make([]*Entity,0)
		} else if table.Rows[val][entity_name] == nil {
			table.Rows[val][entity_name] = make([]*Entity,0)
		}

		table.Rows[val][entity_name] = append(table.Rows[val][entity_name],&entity)
	}

	self.AddStringToSearch(name,&entity)
}

func (self *Database) FetchTable (name string) (*Table) {
	return self.Tables[name]
}

func (self *Database) AddStringToSearch (value string,entity *Entity) {
	result := make(map[string]bool)

	for i:=0;i<len(value);i++ {
		for j:=0;j<len(value);j++ {
			if j+i < len(value) {
				result[value[j:j+i+1]] = true
			}
		}
	}


	var node *SearchTreeNode
	
	if self.SearchTreeHead == nil || self.SearchTreeHead.Children == nil {
		self.SearchTreeHead = &SearchTreeNode{make([]*SearchTreeNode,NUM_OF_CHARS_SUPPORTED),make([]*Entity,0)}
	}

	for key,_ := range(result) {
		node = self.SearchTreeHead.iterateTill(key)
		(*node).Values = append(node.Values,entity)
	}
}

func (self *SearchTreeNode) iterateTill (loc string) (*SearchTreeNode) {
	iterator := &self

	//fmt.Println("\n\n\n\n\n\n\n\n")

	for _,r := range(loc) {
		//fmt.Println(loc,string(r))
		if *iterator== nil || (*iterator).Children == nil {
			*iterator = &SearchTreeNode{make([]*SearchTreeNode,NUM_OF_CHARS_SUPPORTED),make([]*Entity,0)}
		}
		
		if (*iterator).Children[int(r)] == nil || (*iterator).Children[int(r)].Children == nil{
			(*iterator).Children[int(r)] = &SearchTreeNode{make([]*SearchTreeNode,NUM_OF_CHARS_SUPPORTED),make([]*Entity,0)}
		}

		// /fmt.Println(*iterator)
		iterator = &((*iterator).Children[int(r)])
	}
	
	if *iterator == nil {
		*iterator = &SearchTreeNode{make([]*SearchTreeNode,NUM_OF_CHARS_SUPPORTED),make([]*Entity,0)}	
	}

	return *iterator
}

func (self *Database) GetAllTables () (map[string]*Table) {
	return self.Tables
}

func (self *Database) GetTable (name string) (*Table) {
	return self.Tables[name]
}

func (self *Table) QueryRow (rowName , arg string) (map[string][]*Entity) {
	arg = strings.TrimSpace(arg)
	switch arg {
	case MULT_ARGS:
		return self.Rows[rowName]
	default:
		args := strings.Split(arg,QUERY_DELIM)
		ret := make(map[string][]*Entity)
		for _,arg := range(args) {
			ret[arg] = self.Rows[rowName][arg]
		}
		return ret
	}
}

func (self *Table) QueryColumn (colName , arg string) (map[string]*Entity) {
	arg = strings.TrimSpace(arg)
	switch arg {
	case MULT_ARGS:
		t := make(map[string]*Entity)
		t[colName] = self.Columns[colName]
		return t
	default:
		args := strings.Split(arg,QUERY_DELIM)
		ret := make(map[string]*Entity)
		for _,arg := range(args) {
			ret[arg] = self.Columns[colName]
		}
		return ret
	}
}

func (self *Database) SearchFor (str string) ([]*Entity) {
	loc := self.SearchTreeHead.iterateTill(str)
	return loc.Values
}

func (self *Entity) Link (arg *Entity) {
	self.Relations = append(self.Relations,arg)
}

func (self *Entity) Bind (arg *Entity) {
	self.Relations = append(self.Relations,arg)
	arg.Relations = append(arg.Relations,self)
}

func (self *Entity) LinkedTo (arg *Entity) (bool) {
	for _,elem := range(arg.Relations) {
		if self.Name == elem.Name && len(self.Data) == len(elem.Data) {
			return true
		}
	}

	return false
}

func (self *Entity) BindedTo (arg *Entity) (bool) {
	for _,elem := range(arg.Relations) {
		if self.Name == elem.Name && len(self.Data) == len(elem.Data) {
			for _,el := range(self.Relations) {
				if el.Name == arg.Name && len(el.Data) == len(arg.Data) {
					return true
				}
			}
		}
	}

	return false
}

func (arg *Entity) DestroyLinkWith (self *Entity) (bool) {
	for itr,elem := range(arg.Relations) {
		if self.Name == elem.Name && len(self.Data) == len(elem.Data) {
			self.Relations = append(self.Relations[:itr], self.Relations[itr+1:]...)
			return true
		}
	}
	return false
}

func (self *Entity) DestroyBindWith (arg *Entity) (bool) {
	for i,elem := range(arg.Relations) {
		if self.Name == elem.Name && len(self.Data) == len(elem.Data) {
			for j,el := range(self.Relations) {
				if el.Name == arg.Name && len(el.Data) == len(arg.Data) {
					self.Relations = append(self.Relations[:j], self.Relations[j+1:]...)
					arg.Relations = append(arg.Relations[:i], arg.Relations[i+1:]...)
					
					return true
				}
			}
		}
	}

	return false
}

func (self *Entity) GetCommon (args...*Entity) (res []*Entity) {
	if len(args) == 0 {
		
		return
	}
	for _,l1 := range(self.Relations) {
		for _,arg := range(args) {
			if l1.Name == arg.Name && len(l1.Data) == len(arg.Data) {
				res = append(res,l1)
			}
		}
	}
	return res
}