package balus

type LDAP struct {
	ObjectClass                []string `json:"objectClass"`
	UID                        string   `json:"uid"`
	ShadowMin                  int      `json:"shadowMin"`
	RadiusTunnelMediumType     []string `json:"radiusTunnelMediumType"`
	UIDNumber                  int      `json:"uidNumber"`
	ShadowMax                  int      `json:"shadowMax"`
	ShadowLastChange           int      `json:"shadowLastChange"`
	ShadowExpire               int      `json:"shadowExpire"`
	Cn                         []string `json:"cn"`
	RadiusTunnelType           []string `json:"radiusTunnelType"`
	Mail                       []string `json:"mail"`
	RadiusTunnelPrivateGroupID []string `json:"radiusTunnelPrivateGroupId"`
	LdapFirstName              string   `json:"ldapFirstName"`
	LoginShell                 string   `json:"loginShell"`
	GidNumber                  int      `json:"gidNumber"`
	SambaNTPassword            string   `json:"sambaNTPassword"`
	ShadowWarning              int      `json:"shadowWarning"`
	ShadowInactive             int      `json:"shadowInactive"`
	SambaSID                   string   `json:"sambaSID"`
	Gecos                      string   `json:"gecos"`
	Sn                         []string `json:"sn"`
	LdapLastName               string   `json:"ldapLastName"`
	HomeDirectory              string   `json:"homeDirectory"`
	Ou                         []string `json:"ou"`
	Ldapedyid                  []string `json:"ldapedyid"`
}
