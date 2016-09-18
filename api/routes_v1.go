package api

func init() {
	routes["v1"] = routeDefinition{
		"ListNodes":                route{P: "/nodes/", M: "GET", h: listNodes},
		"GetNode":                  route{P: "/nodes/{node}", M: "GET", h: getNode},
		"AddNode":                  route{P: "/nodes/{node}", M: "POST", h: addNode},
		"DelNode":                  route{P: "/nodes/{node}", M: "DELETE", h: delNode},
		"AddNodeVars":              route{P: "/nodes/{node}/vars", M: "POST", h: addNodeVars},
		"ReplaceNodeVars":          route{P: "/nodes/{node}/vars", M: "PUT", h: replaceNodeVars},
		"TriggerNodeVarsProviders": route{P: "/nodes/{node}/vars/providers", M: "PUT", h: triggerNodeVarsProviders},
		"GetNodeVars":              route{P: "/nodes/{node}/vars", M: "GET", h: getNodeVars},
		"GetNodeVar":               route{P: "/nodes/{node}/vars/{var}", M: "GET", h: getNodeVar},
		"LinkNodeToRole":           route{P: "/nodes/{node}/roles/{role}", M: "POST", h: linkNodeToRole},
		"UnlinkNodeFromRole":       route{P: "/nodes/{node}/roles/{role}", M: "DELETE", h: unlinkNodeFromRole},

		"ListRoles":     route{P: "/roles/", M: "GET", h: listRoles},
		"GetRole":       route{P: "/roles/{role}", M: "GET", h: getRole},
		"AddRole":       route{P: "/roles/{role}", M: "POST", h: addRole},
		"DelRole":       route{P: "/roles/{role}", M: "DELETE", h: delRole},
		"AddRoleVars":   route{P: "/roles/{role}/vars", M: "POST", h: addRoleVars},
		"AddRoleParent": route{P: "/roles/{role}/parent/{parent}", M: "POST", h: addRoleParent},
		"DelRoleParent": route{P: "/roles/{role}/parent", M: "DELETE", h: delRoleParent},

		"SiteMap":              route{P: "/system/sytemap", M: "GET", h: sitemapV1},
		"WhoAmI":               route{P: "/system/whoami", M: "GET", h: whoAmI},
		"GetNodeVarsProviders": route{P: "/system/nodevarsproviders", M: "GET", h: getNodeVarsProviders},
		"GetConfig":            route{P: "/system/config", M: "GET", h: getConfig},
	}
}
