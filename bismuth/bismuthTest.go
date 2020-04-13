package main

import (
	"bismuth"
	"fmt"
)

type (
	strct struct {
		f1,f2,f3,f4,f5 string
	}
)

var (
	session *bismuth.Session
	r *bismuth.User
	u1,u2,u3 *bismuth.User
	s1,s2,s3 *bismuth.User
	d1,d2,d3 *bismuth.Database
)

func main () {
	checkCreationFuncs()
}

func checkCreationFuncs () {
	session = bismuth.CreateSession()
	fmt.Println("Session created:",*session)

	r = session.CreateRootUser("super-user","password")
	fmt.Println("super user created:",*r)
	
	u1,u2,u3 = r.CreateUser("jojo","hoho"),r.CreateUser("tralalala","oba"),r.CreateUser("monk","key")
	fmt.Println("Users created:",*u1,*u2,*u3)

	s1,s2,s3 = u1.CreateSegment("san"),u2.CreateSegment("e"),u3.CreateSegment("tux")
	fmt.Println("Segments Created:",*s1,*s2,*s3)

	r.CreateDatabase("s")
	u1.CreateDatabase("ome")
	s1.CreateDatabase("thing")
	fmt.Println("Users Changed:",r.Databases,u1.Databases,s1.Databases,"\n")

	r.CreateTableIn("s","someTable","f1","f2","f3")
	fmt.Println("Table in super-user:",r.Databases["s"].Tables)
	fmt.Println("Table struct:",r.Databases["s"].Tables["someTable"])

	r.GetDatabase("s").FetchTable("someTable").AddRows("f3","f4")
	fmt.Println("The Tables have turned :",r.Databases["s"].Tables["someTable"])

	r.GetDatabase("s").AddEntityIn("someTable","obasan",strct{"f1","f2","f3","f4","f5"})
	r.GetDatabase("s").AddEntityIn("someTable","rarn",strct{"f1","f2","f3","f4","f5"})
	fmt.Println("\nLooks like we got an entity :",r.GetDatabase("s").SearchFor("oba"))

	a := r.GetDatabase("s").SearchFor("oba")
	b := r.GetDatabase("s").SearchFor("rarn")

	a[0].Bind(b[0])

	h := a[0].GetCommon(b...)
	fmt.Println(*h[0])
}