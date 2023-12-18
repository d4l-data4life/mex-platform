package auth

import (
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/known/securitypb"
)

const (
	RoleProducer = "producer"
	RoleConsumer = "consumer"
)

const (
	ResourceIndex  = "index"
	ResourceItems  = "items"
	ResourceConfig = "config"
	ResourceJobs   = "jobs"
	ResourceBlobs  = "blobs"
	ResourceStatus = "status"
	ResourceNotify = "notify"
)

const (
	VerbCreate = "create"
	VerbRead   = "read"
	VerbUpdate = "update"
	VerbDelete = "delete"
	VerbQuery  = "query"
	VerbSend   = "send"
)

type PrivMask = uint64

type PrivMgr struct {
	roles      map[string]*securitypb.Role
	privileges []*securitypb.Privilege
}

func NewPrivMgr() *PrivMgr {
	mgr := PrivMgr{
		roles: map[string]*securitypb.Role{},
	}

	mgr.privileges = []*securitypb.Privilege{
		{Resource: ResourceItems, Verb: VerbCreate},
		{Resource: ResourceItems, Verb: VerbRead},
		{Resource: ResourceItems, Verb: VerbDelete},

		{Resource: ResourceIndex, Verb: VerbCreate},
		{Resource: ResourceIndex, Verb: VerbUpdate},
		{Resource: ResourceIndex, Verb: VerbQuery},
		{Resource: ResourceIndex, Verb: VerbDelete},

		{Resource: ResourceConfig, Verb: VerbRead},
		{Resource: ResourceConfig, Verb: VerbUpdate},

		{Resource: ResourceJobs, Verb: VerbCreate},
		{Resource: ResourceJobs, Verb: VerbRead},

		{Resource: ResourceBlobs, Verb: VerbCreate},
		{Resource: ResourceBlobs, Verb: VerbRead},
		{Resource: ResourceBlobs, Verb: VerbUpdate},
		{Resource: ResourceBlobs, Verb: VerbDelete},

		{Resource: ResourceStatus, Verb: VerbRead},
		{Resource: ResourceNotify, Verb: VerbSend},
	}

	// Assign each privilege its bit mask based on the bit index.
	for b, priv := range mgr.privileges {
		priv.Mask = 1 << b
	}

	mgr.roles[RoleConsumer] = &securitypb.Role{
		Name:        RoleConsumer,
		Description: "Read-only consumer role",
		Mask: 0 |
			mgr.MustPrivMask(ResourceItems, VerbRead) |
			mgr.MustPrivMask(ResourceIndex, VerbQuery) |
			mgr.MustPrivMask(ResourceConfig, VerbRead) |
			mgr.MustPrivMask(ResourceNotify, VerbSend),
	}

	mgr.roles[RoleProducer] = &securitypb.Role{
		Name:        RoleProducer,
		Description: "Read-write producer role",
		Mask: 0 |
			mgr.MustPrivMask(ResourceItems, VerbCreate) |
			mgr.MustPrivMask(ResourceItems, VerbRead) |
			mgr.MustPrivMask(ResourceItems, VerbDelete) |

			mgr.MustPrivMask(ResourceIndex, VerbCreate) |
			mgr.MustPrivMask(ResourceIndex, VerbUpdate) |
			mgr.MustPrivMask(ResourceIndex, VerbQuery) |
			mgr.MustPrivMask(ResourceIndex, VerbDelete) |

			mgr.MustPrivMask(ResourceConfig, VerbRead) |
			mgr.MustPrivMask(ResourceConfig, VerbUpdate) |

			mgr.MustPrivMask(ResourceJobs, VerbRead) |
			mgr.MustPrivMask(ResourceJobs, VerbCreate) |

			mgr.MustPrivMask(ResourceBlobs, VerbCreate) |
			mgr.MustPrivMask(ResourceBlobs, VerbRead) |
			mgr.MustPrivMask(ResourceBlobs, VerbUpdate) |
			mgr.MustPrivMask(ResourceBlobs, VerbDelete) |

			mgr.MustPrivMask(ResourceStatus, VerbRead) |
			mgr.MustPrivMask(ResourceNotify, VerbSend),
	}

	return &mgr
}

func (mgr *PrivMgr) MustPrivMask(resource string, verb string) uint64 {
	for _, priv := range mgr.privileges {
		if priv.Resource == resource && priv.Verb == verb {
			return priv.Mask
		}
	}

	panic("unknown privilege: " + resource + "/" + verb)
}

func (mgr *PrivMgr) ResolveRoles(roleNames []string) (PrivMask, error) {
	var mask uint64

	for _, roleName := range roleNames {
		role, ok := mgr.roles[roleName]
		if !ok {
			return 0, fmt.Errorf("role not found: %s", roleName)
		}
		mask |= role.Mask
	}

	return mask, nil
}
