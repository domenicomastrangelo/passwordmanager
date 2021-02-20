To save a new password:

go run . -action=add -element=github.com

```
Please enter your username: username
Please enter your password: password
Please enter your element: username|password
```

The separator between username and password can be whatever you choose.

To get a saved password:

go run . -action=get -element=github.com

```
Please enter your username: username
Please enter your password: password

Element is: username|password

Element is: anotherusername|anotherpassword
```
