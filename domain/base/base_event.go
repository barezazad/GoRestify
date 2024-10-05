package base

import "GoRestify/pkg/pkg_types"

// types for base domain
const (
	UpdateSetting pkg_types.Event = "setting-update"
	ViewSetting   pkg_types.Event = "setting-view"

	CreateCity pkg_types.Event = "city-create"
	UpdateCity pkg_types.Event = "city-update"
	DeleteCity pkg_types.Event = "city-delete"
	ListCity   pkg_types.Event = "city-list"
	ViewCity   pkg_types.Event = "city-view"

	CreateRegion pkg_types.Event = "region-create"
	UpdateRegion pkg_types.Event = "region-update"
	DeleteRegion pkg_types.Event = "region-delete"
	ListRegion   pkg_types.Event = "region-list"
	ViewRegion   pkg_types.Event = "region-view"

	CreateAccount pkg_types.Event = "account-create"
	UpdateAccount pkg_types.Event = "account-update"
	DeleteAccount pkg_types.Event = "account-delete"
	ListAccount   pkg_types.Event = "account-list"
	ViewAccount   pkg_types.Event = "account-view"

	CreateRole pkg_types.Event = "role-create"
	UpdateRole pkg_types.Event = "role-update"
	DeleteRole pkg_types.Event = "role-delete"
	ListRole   pkg_types.Event = "role-list"
	ViewRole   pkg_types.Event = "role-view"

	CreateUser pkg_types.Event = "user-create"
	UpdateUser pkg_types.Event = "user-update"
	DeleteUser pkg_types.Event = "user-delete"
	ListUser   pkg_types.Event = "user-list"
	ViewUser   pkg_types.Event = "user-view"

	ClearCache     pkg_types.Event = "clear-cache"
	ClearCacheUser pkg_types.Event = "clear-cache-user"

	CreateDocument pkg_types.Event = "document-create"
	UpdateDocument pkg_types.Event = "document-update"
	DeleteDocument pkg_types.Event = "document-delete"
	ListDocument   pkg_types.Event = "document-list"
	ViewDocument   pkg_types.Event = "document-view"
)
