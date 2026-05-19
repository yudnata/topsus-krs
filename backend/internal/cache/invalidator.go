package cache

import "context"

// Invalidator menghapus cache saat data berubah (dipanggil dari service Slice 4+).
type Invalidator struct {
	c *Client
}

func NewInvalidator(c *Client) *Invalidator {
	return &Invalidator{c: c}
}

func (i *Invalidator) OnDosenWrite(ctx context.Context, dosenID string) {
	_ = i.c.Del(ctx,
		KeyDosenList(),
		KeyDosenDetail(dosenID),
		KeyKrsPendingDosen(dosenID),
	)
}

func (i *Invalidator) OnMahasiswaWrite(ctx context.Context, mahasiswaID string, dosenPembimbingID string) {
	keys := []string{
		KeyMahasiswaList(),
		KeyMahasiswaDetail(mahasiswaID),
		KeyKrsMahasiswaList(mahasiswaID),
	}
	if dosenPembimbingID != "" {
		keys = append(keys, KeyKrsPendingDosen(dosenPembimbingID))
	}
	_ = i.c.Del(ctx, keys...)
}

func (i *Invalidator) OnKrsMutation(ctx context.Context, mahasiswaID, dosenPembimbingID string) {
	keys := []string{
		KeyKrsPendingStaff(),
		KeyKrsMahasiswaList(mahasiswaID),
	}
	if dosenPembimbingID != "" {
		keys = append(keys, KeyKrsPendingDosen(dosenPembimbingID))
	}
	_ = i.c.Del(ctx, keys...)
}

func (i *Invalidator) OnDpaChange(ctx context.Context, oldDosenID, newDosenID, mahasiswaID string) {
	i.OnMahasiswaWrite(ctx, mahasiswaID, newDosenID)
	if oldDosenID != "" && oldDosenID != newDosenID {
		_ = i.c.Del(ctx, KeyKrsPendingDosen(oldDosenID))
	}
}
