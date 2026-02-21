package utils

import (
	"os"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

// InitSupabase 初始化 Supabase 客户端
func InitSupabase() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		panic("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	var err error
	Client, err = supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		panic("Failed to initialize Supabase client: " + err.Error())
	}
}
