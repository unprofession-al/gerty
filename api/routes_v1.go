package api

var apiv1Routes = Routes{
	"ListNodes":       Route{p: "/nodes/", m: "GET", h: listNodes},
	"GetNode":         Route{p: "/nodes/{node}", m: "GET", h: getNode},
	"AddNode":         Route{p: "/nodes/{node}", m: "POST", h: addNode},
	"DelNode":         Route{p: "/nodes/{node}", m: "DELETE"},
	"GetNodeVars":     Route{p: "/nodes/{node}/vars", m: "GET", h: getNodeVars},
	"LinkNodeToRole":  Route{p: "/nodes/{node}/roles/{role}", m: "POST", h: linkNodeToRole},
	"ListRoles":       Route{p: "/roles/", m: "GET", h: listRoles},
	"GetRole":         Route{p: "/roles/{role}", m: "GET", h: getRole},
	"AddRole":         Route{p: "/roles/{role}", m: "POST", h: addRole},
	"DelRole":         Route{p: "/roles/{role}", m: "DELETE"},
	"GetRoleParent":   Route{p: "/roles/{role}/parent", m: "GET", h: getRoleParent},
	"AddRoleParent":   Route{p: "/roles/{role}/parent/{parent}", m: "POST", h: addRoleParent},
	"DelRoleParent":   Route{p: "/roles/{role}/parent", m: "DELETE"},
	"GetRoleChildren": Route{p: "/roles/{role}/children", m: "GET"},
	"AddRoleChildren": Route{p: "/roles/{role}/children", m: "POST"},
	"DelRoleChildren": Route{p: "/roles/{role}/children", m: "DELETE"},
	"GetRoleChild":    Route{p: "/roles/{role}/children/{child}", m: "GET"},
	"AddRoleChild":    Route{p: "/roles/{role}/children/{child}", m: "POST"},
	"DelRoleChild":    Route{p: "/roles/{role}/children/{child}", m: "DELETE"},
}
