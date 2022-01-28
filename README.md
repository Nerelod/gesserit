# Gesserit: A simple golang reverse shell manager

## Install

simply git clone this then '''go run gesserit.go''' will do the trick. Compile as you desire.

## Usage
When it starts, it will listen for a session. Once a session is made, it will attempt to
upgrade the shell using a python one liner
> python -c 'import pty; pty.spawn("/bin/bash")'

Once in a session, gesserit commands can be run by prepending 'gesserit' 
Here are the commands
'''
gesserit switch
gesserit list
gesserit quit
'''
That should be self explanatory. 
Switch switches the current session
List lists the sessions
Quit will quit the tool from running.

