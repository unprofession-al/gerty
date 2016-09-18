package api

func init() {
	routes["v1"] = leafs{
		"nodes": leaf{
			E: endpoints{
				"GET": endpoint{N: "ListNodes", h: listNodes},
			},
			L: leafs{
				"{node}": leaf{
					E: endpoints{
						"GET":    endpoint{N: "GetNode", h: getNode},
						"POST":   endpoint{N: "AddNode", h: addNode},
						"DELETE": endpoint{N: "DelNode", h: delNode},
					},
					L: leafs{
						"vars": leaf{
							E: endpoints{
								"GET":  endpoint{N: "GetNodeVars", h: getNodeVars},
								"POST": endpoint{N: "AddeNodeVars", h: addNodeVars},
								"PUT":  endpoint{N: "ReplaceNodesVars", h: replaceNodeVars},
							},
							L: leafs{
								"{var}": leaf{
									E: endpoints{
										"GET": endpoint{N: "GetNodeVar", h: getNodeVar},
									},
								},
								"providers": leaf{
									E: endpoints{
										"PUT": endpoint{N: "TriggerNodeVarsProviders", h: triggerNodeVarsProviders},
									},
								},
							},
						},
						"roles": leaf{
							L: leafs{
								"{role}": leaf{
									E: endpoints{
										"POST":   endpoint{N: "LinkNodeToRole", h: linkNodeToRole},
										"DELETE": endpoint{N: "UnlinkNodeFromRole", h: unlinkNodeFromRole},
									},
								},
							},
						},
					},
				},
			},
		},
		"roles": leaf{
			E: endpoints{
				"GET": endpoint{N: "ListRoles", h: listRoles},
			},
			L: leafs{
				"{role}": leaf{
					E: endpoints{
						"GET":    endpoint{N: "GetRole", h: getRole},
						"POST":   endpoint{N: "AddRole", h: addRole},
						"DELETE": endpoint{N: "DelRole", h: delRole},
					},
					L: leafs{
						"vars": leaf{
							E: endpoints{
								"POST": endpoint{N: "AddRoleVars", h: addRoleVars},
							},
						},
						"parent": leaf{
							E: endpoints{
								"DELETE": endpoint{N: "DelRoleParent", h: delRoleParent},
							},
							L: leafs{
								"{parent}": leaf{
									E: endpoints{
										"POST": endpoint{N: "AddRoleParent", h: addRoleParent},
									},
								},
							},
						},
					},
				},
			},
		},
		"system": leaf{
			L: leafs{
				"sitemap": leaf{
					E: endpoints{
						"GET": endpoint{N: "SiteMap", h: sitemapV1},
					},
				},
				"whoami": leaf{
					E: endpoints{
						"GET": endpoint{N: "WhoAmI", h: whoAmI},
					},
				},
				"nodevarsproviders": leaf{
					E: endpoints{
						"GET": endpoint{N: "GetNodeVarsProviders", h: getNodeVarsProviders},
					},
				},
				"config": leaf{
					E: endpoints{
						"GET": endpoint{N: "GetConfig", h: getConfig},
					},
				},
			},
		},
	}
}
