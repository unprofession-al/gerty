package api

var apiv1Routes = routes{
	"ListNodes":       route{p: "/nodes/", m: "GET", h: listNodes},
	"GetNode":         route{p: "/nodes/{node}", m: "GET", h: getNode},
	"AddNode":         route{p: "/nodes/{node}", m: "POST", h: addNode},
	"DelNode":         route{p: "/nodes/{node}", m: "DELETE"},
	"GetNodeVars":     route{p: "/nodes/{node}/vars", m: "GET", h: getNodeVars},
	"LinkNodeToRole":  route{p: "/nodes/{node}/roles/{role}", m: "POST", h: linkNodeToRole},
	"ListRoles":       route{p: "/roles/", m: "GET", h: listRoles},
	"GetRole":         route{p: "/roles/{role}", m: "GET", h: getRole},
	"AddRole":         route{p: "/roles/{role}", m: "POST", h: addRole},
	"DelRole":         route{p: "/roles/{role}", m: "DELETE"},
	"GetRoleParent":   route{p: "/roles/{role}/parent", m: "GET", h: getRoleParent},
	"AddRoleParent":   route{p: "/roles/{role}/parent/{parent}", m: "POST", h: addRoleParent},
	"DelRoleParent":   route{p: "/roles/{role}/parent", m: "DELETE"},
	"GetRoleChildren": route{p: "/roles/{role}/children", m: "GET"},
	"AddRoleChildren": route{p: "/roles/{role}/children", m: "POST"},
	"DelRoleChildren": route{p: "/roles/{role}/children", m: "DELETE"},
	"GetRoleChild":    route{p: "/roles/{role}/children/{child}", m: "GET"},
	"AddRoleChild":    route{p: "/roles/{role}/children/{child}", m: "POST"},
	"DelRoleChild":    route{p: "/roles/{role}/children/{child}", m: "DELETE"},
}
