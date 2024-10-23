package service

import (
	"fmt"
	"time"
	"net/http"
	"github.com/go-playground/validator/v10"	
	ut "github.com/go-playground/universal-translator"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/mapper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/repository"
)

type PoinService interface {
	SubmissionPoinGet() ([]dto.Pembelian, map[string]interface{}, int)
	SubmissionPoinSuccess(data map[string]interface{}, response map[string]interface{}) (int, error)
	MengaksesArrayInMapAndValidateStruct(data map[string]interface{}, validator *validator.Validate, trans ut.Translator, response map[string]interface{}) ([]dto.Pembelian_Per_Item, int)
	MapToStructSubmissionPoinCancel(pembelian map[string]interface{}) (*dto.KeteranganNotaCancel, *dto.Pembelian)
	SubmissionPoinCancel(data map[string]interface{}, response map[string]interface{}) int
}

type poinService struct {
	repo repository.PoinRepository
} 

func NewServicePoin(repo repository.PoinRepository) PoinService {
	return &poinService{repo: repo}
}

func (ps *poinService) SubmissionPoinGet() ([]dto.Pembelian, map[string]interface{}, int) {
	response := make(map[string]interface{})

	pembelians, err := ps.repo.SubmissionPoinGet()
	if err != nil {
		response["message"] = err.Error()
		return nil, response, http.StatusInternalServerError
	}

	var pembeliansDTO []dto.Pembelian
	for _, data := range pembelians {
		Datadto := mapper.ToPembelianDTO(&data)
		pembeliansDTO = append(pembeliansDTO, Datadto)
	}

	return pembeliansDTO, nil, http.StatusOK
}

func (ps *poinService) SubmissionPoinSuccess(data map[string]interface{}, response map[string]interface{}) (int, error) {
	// inialisasikan validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	trans := helper.TranslatorIDN()
	helper.RegisterCustomValidations(validate, trans)

	PembelianPerItem, statusCode := ps.MengaksesArrayInMapAndValidateStruct(data, validate, trans, response)
	if statusCode != 200 {
		return statusCode, nil
	}

	// menghitung total keuntungan pembelian berdasarkan total keuntungan yang ada pada pembelian per item
	var TotalKeuntunganPembelian float64
	for _, perItemBarang := range PembelianPerItem {
		TotalKeuntunganPembelian = TotalKeuntunganPembelian + perItemBarang.Total_Keuntungan
	}

	layout := "2006-01-02T15:04:05Z07:00" // Sesuaikan format ini dengan format yang ada pada data
	createdAtStr := data["CreatedAt"].(string)
	createdAt, _ := time.Parse(layout, createdAtStr)

	pembelian := &dto.Pembelian{
		ID: uint(data["ID"].(float64)),
		UserID : uint(data["userid"].(float64)),
		Tanggal_Pembelian: data["tanggal_pembelian"].(string),
		Total_Harga: data["total_harga"].(float64),
		CreatedAt: createdAt,
		Total_Keuntungan: TotalKeuntunganPembelian,
		Image: data["image"].(string),
		Status: "success",
	}

	if err := ValidateStructPembelian(pembelian, validate, trans, response); err != nil {
		return http.StatusBadRequest, err
	}

	user, err := ps.repo.GetPoinUser(pembelian.UserID)
	if err != nil {
		response["message"] = "Error Cant Get Poin User: " + err.Error()
		return http.StatusInternalServerError, err
	}

	pembelianDomain := mapper.ToPembelianDomain(pembelian)
	
	var errs error
	tx := ps.repo.BeginNewRepository()
	defer func() {
		if r := recover(); r != nil {
			errs = fmt.Errorf("application panic: %v", r)
			tx.Rollback()
			panic(r)
		} else if errs != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// update field poin in table user
	// kalkuliasi keuntungan
	// poin yang di ambil 7% dari total Keuntungan
	// 1 poin sama dengan 100 rupiah
	keuntunganForPoin := pembelianDomain.Total_Keuntungan * 0.07 
	poin := keuntunganForPoin / 100

	amountPoin := user.Poin + poin

	if errs = ps.repo.UpdatePoinUser(amountPoin, pembelian.UserID, tx); errs != nil {
		response["message"] = "Error Cant Update Poin User:" + errs.Error()
		return http.StatusInternalServerError, errs
	}

	// update table pembelian
	if errs = ps.repo.UpdatePembelianPengajuanPoinSucces(&pembelianDomain, tx); errs != nil {
		response["message"] = errs.Error()
		return http.StatusInternalServerError, errs
	}

	// create setiap pembelian perItems
	for _, pembelianperitem := range PembelianPerItem {
		pembelianperitemDomain := mapper.ToPembelianPerItemDomain(&pembelianperitem)

		if errs = ps.repo.InputPembelianPerItem(&pembelianperitemDomain, tx); errs != nil {
			response["message"] = errs.Error()
			return http.StatusInternalServerError, errs
		}
	}

	response["message"] = "Succesfuly Submission Poin"
	return http.StatusOK, nil
}

// use in submission poin, func to access array of map pembelian per item in map pembelian
// memetakan dari map pembelian berupa array interface to array struct pembelian per item
func (ps *poinService) MengaksesArrayInMapAndValidateStruct(data map[string]interface{}, validate *validator.Validate, trans ut.Translator, response map[string]interface{}) ([]dto.Pembelian_Per_Item, int) {
	var PembelianPerItem []dto.Pembelian_Per_Item

	PerItemBarang := data["pembelian_per_item"].([]interface{})

	for _, perItemBarang := range PerItemBarang {
		itemData := perItemBarang.(map[string]interface{})

		var barang = domain.Barang{ID: uint(itemData["barangID"].(float64))}
		if err := ps.repo.GetDataBarang(&barang); err != nil {
			response["message"] = err.Error()
			return nil, http.StatusInternalServerError
		}

		marginPerBarang := barang.Harga_Barang - barang.Harga_Beli
		totalMargin := itemData["jumlah_barang"].(float64) * marginPerBarang
		
		// Pastikan untuk mengonversi tipe data yang benar
		Item := dto.Pembelian_Per_Item{
			PembelianID:  uint(data["ID"].(float64)),
			BarangID: uint(itemData["barangID"].(float64)),
			Jumlah_Barang: itemData["jumlah_barang"].(float64), 
			Total_Harga:   itemData["total_harga"].(float64),
			Total_Keuntungan: totalMargin,
		}

		if err := ValidateStructPembelianPerItem(&Item, validate, trans, response); err != nil {
			return nil, http.StatusBadRequest
		}

		// if err := validate.Struct(Item); err != nil {
		// 	// Check if the error is of type *validator.InvalidValidationError
		// 	if _, ok := err.(*validator.InvalidValidationError); ok {
		// 		// Handle InvalidValidationError
		// 		response["message"] = "Invalid validation error:" + err.Error()
		// 		return nil, http.StatusBadRequest
		// 	}

		// 	errors := err.(validator.ValidationErrors)
		// 	errorsMessage := errors.Translate(trans)
			
		// 	response["message"] = errorsMessage
		// 	return nil, http.StatusBadRequest
		// }

		PembelianPerItem = append(PembelianPerItem, Item)
	}

	return PembelianPerItem, http.StatusOK 
}


func (ps *poinService) SubmissionPoinCancel(data map[string]interface{}, response map[string]interface{}) int {
	KeteranganNotaCancel, Pembelian := ps.MapToStructSubmissionPoinCancel(data)

	PembelianDomain := mapper.ToPembelianDomain(Pembelian)
	KeteranganNotaCancelDomain := mapper.ToKeteranganNotaCancelDomain(KeteranganNotaCancel)

	
	tx := ps.repo.BeginNewRepository()
	var err error
	defer func() {
		if r := recover(); err != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err = ps.repo.InputKeteranganNotaCancel(&KeteranganNotaCancelDomain, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError
	}

	PembelianDomain.KeteranganNotaCancelID = KeteranganNotaCancelDomain.ID
	if err = ps.repo.UpdatePembelianPengajuanPoinSucces(&PembelianDomain, tx); err != nil {
		response["message"] = err.Error()
		return http.StatusInternalServerError
	}

	response["message"] =  "Berhasil Nota Tidak valid"
	return http.StatusOK
}

func (ps *poinService) MapToStructSubmissionPoinCancel(pembelian map[string]interface{}) (*dto.KeteranganNotaCancel, *dto.Pembelian) {
	KeteranganNota := &dto.KeteranganNotaCancel {
		Desc: pembelian["keterangan"].(string),
	}

	Pembelian := &dto.Pembelian{
		ID:	uint(pembelian["ID"].(float64)),
		UserID: uint(pembelian["userid"].(float64)),
		User: dto.User{
			UserName: pembelian["username"].(string),
			Email: pembelian["email"].(string),
		},
		Tanggal_Pembelian:	pembelian["tanggal_pembelian"].(string),
		Total_Harga:	pembelian["total_harga"].(float64),
		Total_Keuntungan: pembelian["total_keuntungan"].(float64),
		Image: pembelian["image"].(string),
		Status: "cancel",
	}

	return KeteranganNota, Pembelian
}
