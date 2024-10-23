package mapper

import (
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

func ToPembelianPerItemDomain(data *dto.Pembelian_Per_Item) domain.Pembelian_Per_Item {
	return domain.Pembelian_Per_Item {
		ID: data.ID,
		PembelianID: data.PembelianID,
		BarangID: data.BarangID,
		Jumlah_Barang: data.Jumlah_Barang,
		Total_Harga: data.Total_Harga,
		Total_Keuntungan: data.Total_Keuntungan,
	}
}