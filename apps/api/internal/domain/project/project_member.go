package project

import "github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"

// ProjectMember はProject集約内のメンバーを表すエンティティ。
//
// 現在は小規模なProjectを前提に、Membershipを独立集約にせずProjectの子にしている。
// この構成には、メンバー重複の防止、Role変更、最後のAdminの維持を、単一集約の
// トランザクションで確実に守れるという利点がある。また、永続化やユースケースも単純になる。
//
// 一方、1 Projectあたりのメンバーが数千〜数万件になる、メンバー更新の同時実行が多い、
// Project本体とMembershipで読み書きの頻度や保存先が異なる、メンバー一覧を毎回ロードする
// コストが問題になる場合は、ProjectMembershipを独立集約へ切り出す。その場合、重複は
// DBの一意制約で、最後のAdminの維持はロックやトランザクションを伴うドメインサービスで
// 保証する必要があり、集約が小さく競合しにくい代わりに整合性管理が複雑になる。
type ProjectMember struct {
	userID user.UserID
	role   ProjectRole
}

func NewProjectMember(userID user.UserID, role ProjectRole) (ProjectMember, error) {
	if userID.IsZero() {
		return ProjectMember{}, ErrInvalidProjectMemberUserID
	}
	if !role.IsValid() {
		return ProjectMember{}, ErrInvalidProjectRole
	}
	return newProjectMember(userID, role), nil
}

func newProjectMember(userID user.UserID, role ProjectRole) ProjectMember {
	return ProjectMember{userID: userID, role: role}
}

func (m ProjectMember) UserID() user.UserID { return m.userID }

func (m ProjectMember) Role() ProjectRole { return m.role }
