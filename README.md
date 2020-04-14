# bismuth
## A bare bones relational database in and only for golang
<br>
### What is it <br>
Bismuth is a very basic and barebones relational database in golang, it consists of the following data structures to aid its efficiency
- Hash Maps
  These are implemented for tables thus providing a O(1) time to access a given entity unlike O(logN) time taken by B-Trees of SQL and other databases but keep in mind even this has disadvantages (growing and shrinking the base array for more keys in the  hash map)
- An UTF-8 Search Tree
  This is something I came up on my own and it does have an insane space complexity but is better then implementing hashmaps as they do not take extra time to grow and shrink which might be huge disadvantage in normal hashmaps (if you see the code you will understand as I am appending a lot of keys to the searching structure)
- Graphs
  The main reason behind building this database was the fact that my previous crude database I made in python3 for daphne did not have enough capabilities to rival the social giants as of now thus I created this relational Database for pure efficiency without one thought about the space complexity, <br>Now using graphs is extremely easy in the package because all you have to do is call the Link or Bind method on an entity
<br><br>
### Why to use it <br>
As I have stated the database had been made with efficiency in mind not space and thus it provides fast accessing of data due to various data structures, thus it is perfect for small and quick projects as there is no need to learn any syntax to use the database, as everything can be done by functions
<br><br>
### Structure <br>
Now I have specially made this column to explain the odd structure of my database, which is somewhat inspired by linux and my odd habit to keep my computer super neat<br><br>
So onto the point, First of all there is a session you need to create a session to actually start using the database then in the session there should always be a root user to create other users, this was done to create segments and break the databases into parts for ease of access, and by the way all the databases of the all the users are accessible by the root and a non-root user can only create a segment which in essense is a dead user with a name and no password (about the password, that is there for ethics and not encryppted so don't rely on it).<br><br>
Now every user can have multiple databases and each database can have several Tables which can have several entities now entities is a packet containing data about the value such as its relations,fields etcetra,etcetra, another fun easter egg in this code is the fact that you can just put a struct as an argument and by using deep reflect not only can it decipher the fields and automaticaly store them but also create a new table in case the given struct can not fit in any table
