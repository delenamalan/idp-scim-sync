package model

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func TestMember_GobEncode(t *testing.T) {
	tests := []struct {
		name    string
		m       Member
		wantErr bool
	}{
		{
			name: "Test Member GobEncode",
			m: Member{
				IPID:     "1",
				SCIMID:   "1",
				Email:    "member.1@mail.com",
				HashCode: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GobEncode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Member.GobEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)
			if err := enc.Encode(tt.m.IPID); err != nil {
				t.Fatal(err)
			}

			if err := enc.Encode(tt.m.Email); err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, b.Bytes()) {
				t.Errorf("Group.GobEncode() = %v\n, want %v\n", got, b.Bytes())
			}
		})
	}
}

func TestMember_SetHashCode(t *testing.T) {
	tests := []struct {
		name   string
		member Member
		want   Member
	}{
		{
			name: "success",
			member: Member{
				IPID:     "1",
				SCIMID:   "1",
				Email:    "user.1@mail.com",
				Status:   "ACTIVE",
				HashCode: "test",
			},
			want: Member{
				IPID:  "1",
				Email: "user.1@mail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.member.SetHashCode()
			tt.want.SetHashCode()

			got := tt.member.HashCode
			if got != tt.want.HashCode {
				t.Errorf("Member.SetHashCode() = %s, want %s", got, tt.want.HashCode)
			}
		})
	}
}

func TestMembersResult_SetHashCode(t *testing.T) {
	m1 := &Member{IPID: "1", SCIMID: "1", Email: "group.1@mail.com", Status: "ACTIVE"}
	m2 := &Member{IPID: "2", SCIMID: "2", Email: "group.2@mail.com", Status: "ACTIVE"}
	m3 := &Member{IPID: "3", SCIMID: "3", Email: "group.3@mail.com", Status: "ACTIVE"}

	m1.SetHashCode()
	m2.SetHashCode()
	m3.SetHashCode()

	mr1 := MembersResult{
		Items:     3,
		Resources: []*Member{m1, m2, m3},
	}
	mr1.SetHashCode()

	mr2 := MembersResult{
		Items:     3,
		Resources: []*Member{m2, m3, m1},
	}
	mr2.SetHashCode()

	mr3 := MembersResult{
		Items:     3,
		Resources: []*Member{m3, m2, m1},
	}
	mr3.SetHashCode()

	t.Logf("mr1: %s\n", utils.ToJSON(mr1))
	t.Logf("mr2: %s\n", utils.ToJSON(mr2))
	t.Logf("mr3: %s\n", utils.ToJSON(mr3))

	t.Logf("mr1.HashCode: %s\n", mr1.HashCode)
	t.Logf("mr2.HashCode: %s\n", mr2.HashCode)
	t.Logf("mr3.HashCode: %s\n", mr3.HashCode)

	if mr1.HashCode != mr2.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if mr1.HashCode != mr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if mr2.HashCode != mr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
}

func TestGroupsMembersResult_SetHashCode(t *testing.T) {
	m1 := &Member{IPID: "1", SCIMID: "1", Email: "group.1@mail.com"}
	m2 := &Member{IPID: "2", SCIMID: "2", Email: "group.2@mail.com"}
	m3 := &Member{IPID: "3", SCIMID: "3", Email: "group.3@mail.com"}

	m1.SetHashCode()
	m2.SetHashCode()
	m3.SetHashCode()

	gm1 := &GroupMembers{Group: &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"}, Resources: []*Member{m1, m2, m3}}
	gm2 := &GroupMembers{Group: &Group{IPID: "2", SCIMID: "2", Name: "group", Email: "group.2@mail.com"}, Resources: []*Member{m2, m1, m3}}
	gm3 := &GroupMembers{Group: &Group{IPID: "3", SCIMID: "3", Name: "group", Email: "group.3@mail.com"}, Resources: []*Member{m1, m3, m2}}

	gm1.SetHashCode()
	gm2.SetHashCode()
	gm3.SetHashCode()

	gmr1 := GroupsMembersResult{
		Items:     3,
		Resources: []*GroupMembers{gm1, gm2, gm3},
	}
	gmr1.SetHashCode()

	gmr2 := GroupsMembersResult{
		Items:     3,
		Resources: []*GroupMembers{gm2, gm3, gm1},
	}
	gmr2.SetHashCode()

	gmr3 := GroupsMembersResult{
		Items:     3,
		Resources: []*GroupMembers{gm3, gm2, gm1},
	}
	gmr3.SetHashCode()

	gmr4 := MergeGroupsMembersResult(&gmr2, &gmr1, &gmr3)
	gmr4.SetHashCode()
	gmr5 := MergeGroupsMembersResult(&gmr3, &gmr2, &gmr1)
	gmr5.SetHashCode()

	t.Logf("gmr4: %s\n", utils.ToJSON(gmr4))
	t.Logf("gmr5: %s\n", utils.ToJSON(gmr5))

	t.Logf("gmr4.HashCode: %s\n", gmr4.HashCode)
	t.Logf("gmr5.HashCode: %s\n", gmr5.HashCode)

	if gmr1.HashCode != gmr2.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if gmr1.HashCode != gmr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if gmr2.HashCode != gmr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}

	if gmr5.HashCode != gmr4.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal: gmr5-> %s, gmr4-> %s", gmr5.HashCode, gmr4.HashCode)
	}
}

func TestGroupMembers_SetHashCode(t *testing.T) {
	tests := []struct {
		name         string
		groupMembers GroupMembers
		want         GroupMembers
	}{
		{
			name: "success",
			groupMembers: GroupMembers{
				Items:    3,
				HashCode: "test",
				Group:    &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"},
				Resources: []*Member{
					{IPID: "1", SCIMID: "1", Email: "group.1@mail.com", Status: "ACTIVE"},
					{IPID: "2", SCIMID: "2", Email: "group.2@mail.com", Status: "ACTIVE"},
					{IPID: "3", SCIMID: "3", Email: "group.3@mail.com", Status: "ACTIVE"},
				},
			},
			want: GroupMembers{
				Items: 3,
				Group: &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"},
				Resources: []*Member{
					{IPID: "3", SCIMID: "3", Email: "group.3@mail.com", Status: "ACTIVE"},
					{IPID: "1", SCIMID: "1", Email: "group.1@mail.com", Status: "ACTIVE"},
					{IPID: "2", SCIMID: "2", Email: "group.2@mail.com", Status: "ACTIVE"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.groupMembers.SetHashCode()
			tt.want.SetHashCode()

			got := tt.groupMembers.HashCode

			if got != tt.want.HashCode {
				t.Errorf("GroupMembers.SetHashCode() = %s, want %s", got, tt.want.HashCode)
				t.Errorf("GroupMembers.SetHashCode() = %+v, want %+v", tt.groupMembers, tt.want)
			}
		})
	}
}

func TestGroupsMembersResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		HashCode  string
		Resources []*GroupMembers
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:   "empty",
			fields: fields{},
			want: []byte(`{
  "items": 0,
  "hashCode": "",
  "resources": []
}`),
			wantErr: false,
		},
		{
			name: "not empty",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: &Group{
							IPID:   "1",
							SCIMID: "1",
							Name:   "group 1",
							Email:  "group.1@mai.com",
						},
					},
				},
			},
			want: []byte(`{
  "items": 1,
  "hashCode": "test",
  "resources": [
    {
      "items": 1,
      "group": {
        "ipid": "1",
        "scimid": "1",
        "name": "group 1",
        "email": "group.1@mai.com",
        "hashCode": ""
      },
      "resources": null
    }
  ]
}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &GroupsMembersResult{
				Items:     tt.fields.Items,
				Resources: tt.fields.Resources,
				HashCode:  tt.fields.HashCode,
			}
			got, err := ur.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersResult.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UsersResult.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}
