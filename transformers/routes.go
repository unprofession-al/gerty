package transformers

func init() {
	routes = routeDefinition{
		"TreeRenderer": route{p: "/tree/data.json", m: "GET", h: treeRenderer},
		"TreeIndex":    route{p: "/tree", m: "GET", h: treeIndex},

		"AnsibleRenderer": route{p: "/ansible", m: "GET", h: ansibleRenderer},
	}
}
