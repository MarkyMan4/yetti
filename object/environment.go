package object

type Environment struct {
	definitions map[string]Object
	parent      *Environment
}

func NewEnvironment() *Environment {
	defs := make(map[string]Object)

	return &Environment{definitions: defs}
}

func CreateChildEnvironment(env *Environment) *Environment {
	defs := make(map[string]Object)
	child := &Environment{definitions: defs, parent: env}

	return child
}

func (e *Environment) Get(ident string) (Object, bool) {
	obj, ok := e.definitions[ident]

	// if variable not found in this environment, recursively check parent environments
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(ident)
	}

	return obj, ok
}

// sets variable with ident for environment
// forceCurrent flag can be used to tell it not to try to set the value in parent
// environments if it is found there. This is used for var statements, function definitions
// and setting argument values during function calls
func (e *Environment) Set(ident string, obj Object, forceCurrent bool) {
	if forceCurrent {
		e.definitions[ident] = obj
		return
	}

	env := e
	_, ok := env.definitions[ident]

	// if ident not found in current environment, check if it exists in parent envrionments
	for !ok && env.parent != nil {
		env = env.parent
		_, ok = env.definitions[ident]
	}

	// set the value in the envrionment it was found in
	// it must exist somewhere, it is on the caller to either forceCurrent or
	// ensure that the variable exists
	env.definitions[ident] = obj
}

func (e *Environment) GetEnvMap() map[string]Object {
	return e.definitions
}

func (e *Environment) GetParentEnv() *Environment {
	return e.parent
}
