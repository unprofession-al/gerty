package transformers

func init() {
	routes = routeDefinition{
		"TreeRenderer": route{p: "/tree", m: "GET", h: treeRenderer},

		"AnsibleRenderer": route{p: "/ansible", m: "GET", h: ansibleRenderer},
	}
}
