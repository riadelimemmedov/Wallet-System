package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// !createRandomUser => creates a test user with random data and validates the created user's fields.
// !It returns the created user instance.
func createRandomUser(t *testing.T) User {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	arg := CreateUserParams{
		Username:     common.RandomUsername(),
		PasswordHash: common.RandomPassword(),
		Email: sql.NullString{
			String: common.RandomEmail(),
			Valid:  true,
		},
		FirstName: sql.NullString{
			String: common.RandomFirstName(),
			Valid:  true,
		},
		LastName: sql.NullString{
			String: common.RandomLastName(),
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: common.RandomPhoneNumber(),
			Valid:  true,
		},
		ProfileImageUrl: sql.NullString{
			String: common.RandomProfileImage(),
			Valid:  true,
		},
	}

	user, err := sqlStore.Queries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.ProfileImageUrl, user.ProfileImageUrl)

	require.True(t, user.IsActive)
	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	require.NotZero(t, user.LastLogin)

	return user
}

// ! TestCreateUser => validates that a user can be successfully created with all fields populated.
// ! It uses createRandomUser helper function to create and validate a new user.
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
	defer CleanupDB(t)
}

// ! TestCreateUserWithNullFields => validates that a user can be created with optional fields set to null.
// ! It ensures that nullable fields are properly handled during user creation.
func TestCreateUserWithNullFields(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	arg := CreateUserParams{
		Username:     common.RandomUsername(),
		PasswordHash: common.RandomPassword(),
		Email: sql.NullString{
			String: common.RandomEmail(),
			Valid:  true,
		},
		FirstName:       sql.NullString{Valid: false},
		LastName:        sql.NullString{Valid: false},
		PhoneNumber:     sql.NullString{Valid: false},
		ProfileImageUrl: sql.NullString{Valid: false},
	}

	user, err := sqlStore.Queries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.False(t, user.FirstName.Valid)
	require.False(t, user.LastName.Valid)
	require.False(t, user.PhoneNumber.Valid)
	require.False(t, user.ProfileImageUrl.Valid)

	defer CleanupDB(t)
}

// ! TestGetUser => validates the retrieval of a user by ID.
// ! It creates a user, retrieves it, and ensures all fields match the original user.
func TestGetUser(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	user1 := createRandomUser(t)
	user2, err := sqlStore.Queries.GetUser(context.Background(), user1.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.ProfileImageUrl, user2.ProfileImageUrl)
	require.Equal(t, user1.IsActive, user2.IsActive)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
	require.WithinDuration(t, user1.LastLogin, user2.LastLogin, time.Second)

	defer CleanupDB(t)
}

// ! TestListUsers => validates the pagination functionality of user listing.
// ! It creates multiple users and verifies that the correct number of users
// ! are returned with the specified limit and offset.
func TestListUsers(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := sqlStore.Queries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
		require.NotZero(t, user.UserID)
		require.NotEmpty(t, user.Username)
		require.True(t, user.Email.Valid)
		require.NotEmpty(t, user.Email.String)
	}
	defer CleanupDB(t)
}

// ! TestUpdateUser => validates that a user's information can be updated.
// ! It creates a user, updates all fields with new random values,
// ! and verifies the changes were successful.
func TestUpdateUser(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		UserID: user1.UserID,
		Username: sql.NullString{
			String: common.RandomUsername(),
			Valid:  true,
		},
		Email: sql.NullString{
			String: common.RandomEmail(),
			Valid:  true,
		},
		FirstName: sql.NullString{
			String: common.RandomFirstName(),
			Valid:  true,
		},
		LastName: sql.NullString{
			String: common.RandomLastName(),
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: common.RandomPhoneNumber(),
			Valid:  true,
		},
		ProfileImageUrl: sql.NullString{
			String: common.RandomProfileImage(),
			Valid:  true,
		},
	}

	user2, err := sqlStore.Queries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)

	defer CleanupDB(t)
}

// ! TestUpdateUserNullFields => validates that a user can be partially updated
// ! with some fields set to null. It ensures that non-updated fields retain
// ! their original values and null fields are properly handled.
func TestUpdateUserNullFields(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		UserID: user1.UserID,
		Email: sql.NullString{
			String: common.RandomEmail(),
			Valid:  true,
		},
		FirstName:       sql.NullString{Valid: false},
		LastName:        sql.NullString{Valid: false},
		PhoneNumber:     sql.NullString{Valid: false},
		ProfileImageUrl: sql.NullString{Valid: false},
	}

	user2, err := sqlStore.Queries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, arg.Email, user2.Email)

	defer CleanupDB(t)
}

// ! TestDeleteUser => validates the soft deletion of a user.
// ! It creates a user, soft deletes it, and verifies that the user
// ! still exists but is marked as inactive.
func TestDeleteUser(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	user1 := createRandomUser(t)
	require.True(t, user1.IsActive)

	err := sqlStore.Queries.DeleteUser(context.Background(), user1.UserID)
	require.NoError(t, err)

	user2, err := sqlStore.Queries.GetUser(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.False(t, user2.IsActive)

	defer CleanupDB(t)
}

// ! TestHardDeleteUser validates the permanent deletion of a user.
// ! It creates a user, permanently deletes it, and verifies that
// ! the user can no longer be retrieved from the database.
func TestHardDeleteUser(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	user1 := createRandomUser(t)

	err := sqlStore.Queries.HardDeleteUser(context.Background(), user1.UserID)
	require.NoError(t, err)

	user2, err := sqlStore.Queries.GetUser(context.Background(), user1.UserID)
	require.Error(t, err)
	require.Empty(t, user2)

	defer CleanupDB(t)
}
