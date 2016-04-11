package api

func init() {
	routes["v1"] = routeDefinition{
		"ListNodes":          route{p: "/nodes/", m: "GET", h: listNodes},
		"GetNode":            route{p: "/nodes/{node}", m: "GET", h: getNode},
		"AddNode":            route{p: "/nodes/{node}", m: "POST", h: addNode},
		"DelNode":            route{p: "/nodes/{node}", m: "DELETE", h: delNode},
		"AddNodeVars":        route{p: "/nodes/{node}/vars", m: "POST", h: addNodeVars},
		"ReplaceNodeVars":    route{p: "/nodes/{node}/vars", m: "PUT", h: replaceNodeVars},
		"GetNodeVars":        route{p: "/nodes/{node}/vars", m: "GET", h: getNodeVars},
		"LinkNodeToRole":     route{p: "/nodes/{node}/roles/{role}", m: "POST", h: linkNodeToRole},
		"UnlinkNodeFromRole": route{p: "/nodes/{node}/roles/{role}", m: "DELETE", h: unlinkNodeFromRole},

		"ListRoles":     route{p: "/roles/", m: "GET", h: listRoles},
		"GetRole":       route{p: "/roles/{role}", m: "GET", h: getRole},
		"AddRole":       route{p: "/roles/{role}", m: "POST", h: addRole},
		"DelRole":       route{p: "/roles/{role}", m: "DELETE", h: delRole},
		"AddRoleVars":   route{p: "/roles/{role}/vars", m: "POST", h: addRoleVars},
		"AddRoleParent": route{p: "/roles/{role}/parent/{parent}", m: "POST", h: addRoleParent},
		"DelRoleParent": route{p: "/roles/{role}/parent", m: "DELETE", h: delRoleParent},
	}
}
