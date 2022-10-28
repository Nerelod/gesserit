# Gesserit: A simple golang reverse shell manager

## Install

simply git clone this then '''go run gesserit.go''' will do the trick. Compile as you desire.

## Usage
When it starts, it will listen for a session. Once a session is made, it will attempt to
upgrade the shell using a python one liner
> python -c 'import pty; pty.spawn("/bin/bash")'

Once in a session, gesserit commands can be run by prepending 'gesserit' 
Here are the commands
```
gesserit switch
gesserit list
gesserit grouplist
gesserit add
gesserit remove
gesserit groupsend
gesserit hush
gesserit yell
gesserit quit
```

### Switch
Switch to another session shell
```
gesserit switch <session id to switch to>
```
### List
List all sessions
```
gesserit list
```
### Group list
list sessions added to the group
```
gesserit grouplist
```
### Add
Add session to the group
```
gesserit add <session id to add>
```
### Remove
Remove session from group
```
gesserit remove <session id to remove>
```
### Group Send
Send a command to all sessions in the group
```
gesserit groupsend <command(s) to send>
```
### Hush
tell gesserit not to announce new shells
```
gesserit hush
```
### Yell
tell gesserit to announce new shells
```
gesserit yell
```
### Quit
quit the tool
````
gesserit quit
```

