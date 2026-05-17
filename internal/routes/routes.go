package routes

import (
	"reflect-backend/internal/config"
	"reflect-backend/internal/database"
	"reflect-backend/internal/middleware"

	"reflect-backend/internal/modules/address"
	"reflect-backend/internal/modules/auth"
	"reflect-backend/internal/modules/cart"
	"reflect-backend/internal/modules/category"
	"reflect-backend/internal/modules/collection"
	"reflect-backend/internal/modules/order"
	"reflect-backend/internal/modules/product"
	"reflect-backend/internal/modules/profile"
	"reflect-backend/internal/modules/upload"
	"reflect-backend/internal/modules/user"
	"reflect-backend/internal/modules/wishlist"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	api := router.Group("/api/v1")

	userRepository := user.NewUserRepository(database.DB)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	authService := auth.NewAuthService(userRepository, cfg)
	authHandler := auth.NewAuthHandler(authService)

	profileService := profile.NewProfileService(userRepository)
	profileHandler := profile.NewProfileHandler(profileService)

	productRepository := product.NewProductRepository(database.DB)
	productService := product.NewProductService(productRepository)
	productHandler := product.NewProductHandler(productService)

	categoryRepository := category.NewCategoryRepository(database.DB)
	categoryService := category.NewCategoryService(categoryRepository)
	categoryHandler := category.NewCategoryHandler(categoryService)

	collectionRepository := collection.NewCollectionRepository(database.DB)
	collectionService := collection.NewCollectionService(collectionRepository)
	collectionHandler := collection.NewCollectionHandler(collectionService)

	cartRepository := cart.NewCartRepository(database.DB)
	cartService := cart.NewCartService(cartRepository)
	cartHandler := cart.NewCartHandler(cartService)

	wishlistRepository := wishlist.NewWishlistRepository(database.DB)
	wishlistService := wishlist.NewWishlistService(wishlistRepository)
	wishlistHandler := wishlist.NewWishlistHandler(wishlistService)

	addressRepository := address.NewAddressRepository(database.DB)
	addressService := address.NewAddressService(addressRepository)
	addressHandler := address.NewAddressHandler(addressService)

	orderRepository := order.NewOrderRepository(database.DB)
	orderService := order.NewOrderService(orderRepository)
	orderHandler := order.NewOrderHandler(orderService)

	cloudinaryClient, _ := upload.NewCloudinary(cfg)

	uploadService := upload.NewUploadService(cloudinaryClient)
	uploadHandler := upload.NewUploadHandler(uploadService)

	authMiddleware := middleware.AuthMiddleware(cfg)

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "pong",
		})
	})

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	productRoutes := api.Group("/products")
	{
		productRoutes.GET("", productHandler.GetAllProducts)
		productRoutes.GET("/new-arrivals", productHandler.GetNewArrivals)
		productRoutes.GET("/featured", productHandler.GetFeaturedProducts)
		productRoutes.GET("/category", productHandler.GetProductsByCategory)
		productRoutes.GET("/slug/:slug", productHandler.GetProductBySlug)
		productRoutes.GET("/:id", productHandler.GetProductByID)
	}

	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.GET("", categoryHandler.GetAllCategories)
		categoryRoutes.GET("/active", categoryHandler.GetActiveCategories)
		categoryRoutes.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
	}

	collectionRoutes := api.Group("/collections")
	{
		collectionRoutes.GET("", collectionHandler.GetAllCollections)
		collectionRoutes.GET("/active", collectionHandler.GetActiveCollections)
		collectionRoutes.GET("/featured", collectionHandler.GetFeaturedCollections)
		collectionRoutes.GET("/slug/:slug", collectionHandler.GetCollectionBySlug)
		collectionRoutes.GET("/:id", collectionHandler.GetCollectionByID)
	}

	protectedRoutes := api.Group("")
	protectedRoutes.Use(authMiddleware)
	{
		userRoutes := protectedRoutes.Group("/users")
		{
			userRoutes.GET("", middleware.RoleMiddleware("admin"), userHandler.GetAllUsers)
			userRoutes.GET("/:id", userHandler.GetUserByID)
			userRoutes.DELETE("/:id", middleware.RoleMiddleware("admin"), userHandler.DeleteUser)
		}

		profileRoutes := protectedRoutes.Group("/me")
		{
			profileRoutes.GET("", profileHandler.GetMyProfile)
			profileRoutes.PUT("", profileHandler.UpdateMyProfile)
		}

		cartRoutes := protectedRoutes.Group("/cart")
		{
			cartRoutes.GET("", cartHandler.GetMyCart)
			cartRoutes.POST("", cartHandler.AddToCart)
			cartRoutes.PUT("/:id", cartHandler.UpdateCartItem)
			cartRoutes.DELETE("/:id", cartHandler.RemoveCartItem)
			cartRoutes.DELETE("", cartHandler.ClearCart)
		}

		wishlistRoutes := protectedRoutes.Group("/wishlist")
		{
			wishlistRoutes.GET("", wishlistHandler.GetMyWishlist)
			wishlistRoutes.POST("", wishlistHandler.AddToWishlist)
			wishlistRoutes.DELETE("/:id", wishlistHandler.RemoveWishlistItem)
			wishlistRoutes.DELETE("/product/:productId", wishlistHandler.RemoveWishlistByProduct)
		}

		addressRoutes := protectedRoutes.Group("/addresses")
		{
			addressRoutes.GET("", addressHandler.GetMyAddresses)
			addressRoutes.POST("", addressHandler.CreateAddress)
			addressRoutes.PUT("/:id", addressHandler.UpdateAddress)
			addressRoutes.DELETE("/:id", addressHandler.DeleteAddress)
		}

		orderRoutes := protectedRoutes.Group("/orders")
		{
			orderRoutes.GET("", orderHandler.GetMyOrders)
			orderRoutes.GET("/:id", orderHandler.GetMyOrderByID)
			orderRoutes.GET("/number/:orderNumber", orderHandler.GetMyOrderByNumber)
		}

		protectedRoutes.POST("/checkout", orderHandler.Checkout)

		uploadRoutes := protectedRoutes.Group("/upload")
		{
			uploadRoutes.POST("", uploadHandler.UploadImage)
		}
	}

	

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(authMiddleware)
	adminRoutes.Use(middleware.RoleMiddleware("admin"))
	{
		adminProducts := adminRoutes.Group("/products")
		{
			adminProducts.POST("", productHandler.CreateProduct)
			adminProducts.PUT("/:id", productHandler.UpdateProduct)
			adminProducts.DELETE("/:id", productHandler.DeleteProduct)
		}

		adminCategories := adminRoutes.Group("/categories")
		{
			adminCategories.POST("", categoryHandler.CreateCategory)
			adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
			adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		adminCollections := adminRoutes.Group("/collections")
		{
			adminCollections.POST("", collectionHandler.CreateCollection)
			adminCollections.PUT("/:id", collectionHandler.UpdateCollection)
			adminCollections.DELETE("/:id", collectionHandler.DeleteCollection)
		}

		adminOrders := adminRoutes.Group("/orders")
		{
			adminOrders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}
	}

	
	
}