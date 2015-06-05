package simredis

type Command struct {
	CommandName string
	Args        []interface{}
}

func NewCommand(commandName string, args ...interface{}) (command Command) {
	command.CommandName = commandName
	command.Args = args
	return command
}

func (this Command) Append(args interface{}) Command {
	this.Args = append(this.Args, args)
	return this
}

type Commonds []Command

func (this Commonds) Append(command Command) Commonds {
	this = append(this, command)
	return this
}
