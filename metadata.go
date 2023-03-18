package shdb

func (m *Metadata) TypeId() TypeId {
	res := TypeId{}
	res.SetType([4]byte(m.Type))
	res.SetUuidBytes(m.Uuid)
	return res
}
