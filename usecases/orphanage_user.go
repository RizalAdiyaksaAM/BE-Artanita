package usecases

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
	"tugas-akhir/drivers/cloudinary"
	dto_base "tugas-akhir/dto/base"
	dto "tugas-akhir/dto/orphanage"
	"tugas-akhir/entities"
	"tugas-akhir/repositories"
	err_util "tugas-akhir/utils/error"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OrphanageUserUsecase interface {
	CreateOrphanageUser(c echo.Context, req *dto.OrphanageUserRequest) error
	GetOrphanageUserAll(
		c echo.Context, req *dto_base.PaginationRequest,
		searchName, filterAddress, filterEducation, filterPosition, filterAge string,
	) (*[]dto.OrphanageUserResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetOrphanageUserByID(c echo.Context, id uuid.UUID) (*dto.OrphanageUserResponse, error)
	UpdateOrphanageUser(c echo.Context, id uuid.UUID, req *dto.OrphanageUserRequest) error
	DeleteOrphanageUser(c echo.Context, id uuid.UUID) error
	GetOrphanageUserByPosition(c echo.Context, position string) (*entities.OrphanageUser, error)
}

type orphanageUserUsecase struct {
	userRepo          repositories.OrphanageUserRepository
	cloudinaryService cloudinary.CloudinaryService
}

func NewOrphanageUserUsecase(userRepo repositories.OrphanageUserRepository, cloudinaryService cloudinary.CloudinaryService) OrphanageUserUsecase {
	return &orphanageUserUsecase{
		userRepo:          userRepo,
		cloudinaryService: cloudinaryService,
	}
}

func (u *orphanageUserUsecase) CreateOrphanageUser(c echo.Context, req *dto.OrphanageUserRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Cek apakah "age" ada dalam form-data dan konversikan ke int
	ageStr := c.FormValue("age") // mengambil nilai "age" dari form-data
	if ageStr != "" {
		// Coba konversi string ke int
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			// Jika konversi gagal, kembalikan error
			return fmt.Errorf("invalid age value: %v", err)
		}
		req.Age = &age // set nilai age ke struct request setelah konversi
	}

	// Menangani file gambar yang dikirim melalui form-data
	formHeader, err := c.FormFile("image")
	if err != nil {
		fmt.Println("error getting form file")
		return err
	}
	formFile, err := formHeader.Open()
	if err != nil {
		fmt.Println("error opening form file")
		return err
	}
	defer formFile.Close()

	// Upload gambar ke Cloudinary
	imageURL, err := u.cloudinaryService.UploadImage(ctx, formFile, "artanita/orphanage/user")
	if err != nil {
		return err
	}

	// Membuat pointer untuk URL gambar yang telah diupload
	imageURLPtr := &imageURL

	// Membuat objek user berdasarkan data yang diterima
	user := &entities.OrphanageUser{
		ID:        uuid.New(),
		Name:      req.Name,
		Address:   req.Address,
		Image:     imageURLPtr,
		Age:       req.Age, // Pastikan age sudah diubah menjadi int
		Position:  req.Position,
		Education: req.Education,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	// Menyimpan user ke dalam database
	err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *orphanageUserUsecase) GetOrphanageUserAll(
	c echo.Context, req *dto_base.PaginationRequest,
	searchName, filterAddress, filterEducation, filterPosition, filterAge string,
) (*[]dto.OrphanageUserResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next string
		prev string
	)

	// Pastikan page tidak lebih kecil dari 1
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Page > 1 {
		prev = baseURL + strconv.Itoa(req.Page-1)
	}

	// Mendapatkan data pengguna dengan filter dan search
	users, totalData, err := u.userRepo.GetUserAll(
		ctx, req, searchName, filterAddress, filterEducation, filterPosition, filterAge,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// Membentuk response dari data yang diperoleh
	var responses []dto.OrphanageUserResponse
	for _, u := range users {
		responses = append(responses, dto.OrphanageUserResponse{
			ID:        u.ID.String(),
			Name:      u.Name,
			Address:   u.Address,
			Age:       u.Age,
			Image:     u.Image,
			Position:  u.Position,
			Education: u.Education,
		})
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	// Menangani kasus ketika halaman yang diminta lebih besar dari total halaman
	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	} else {
		next = baseURL + strconv.Itoa(req.Page+1)
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return &responses, paginationMetadata, link, nil
}

func (u *orphanageUserUsecase) GetOrphanageUserByID(c echo.Context, id uuid.UUID) (*dto.OrphanageUserResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := &dto.OrphanageUserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Address:   user.Address,
		Age:       user.Age,
		Image:     user.Image,
		Position:  user.Position,
		Education: user.Education,
	}

	return response, nil
}

func (u *orphanageUserUsecase) UpdateOrphanageUser(c echo.Context, id uuid.UUID, req *dto.OrphanageUserRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Address = req.Address
	user.Position = req.Position
	user.Education = req.Education

	err = u.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *orphanageUserUsecase) DeleteOrphanageUser(c echo.Context, id uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *orphanageUserUsecase) GetOrphanageUserByPosition(c echo.Context, position string) (*entities.OrphanageUser, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := u.userRepo.GetUserByPosition(ctx, position)
	if err != nil {
		return nil, err
	}

	return user, nil
}
