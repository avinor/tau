package config

// Dependency towards another tau deployment. Source can either be a relative / absolute path
// (start with . or / in that case) to a file or a directory.
//
// For each dependency it will create a remote_state data source to retrieve the values from
// dependency. Backend configuration will be read from the dependency file. To override attributes
// define the backend block in dependency and only define the attributes that should be overridden.
// For instance it can be useful to override token attribute if current module and dependency module
// use different token's for authentication
type Dependency struct {
	Name   string `hcl:"name,label"`
	Source string `hcl:"source,attr"`

	Backend *Backend `hcl:"backend,block"`
}
