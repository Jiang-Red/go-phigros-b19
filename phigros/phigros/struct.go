package phigros

import "time"

type PhigrosStruct interface {
	Settings | User | Summary
}
type Settings struct {
	ChordSupport      bool    `json:"chordSupport"`
	FcAPIndicator     bool    `json:"fcAPIndicator"`
	EnableHitSound    bool    `json:"enableHitSound"`
	LowResolutionMode bool    `json:"lowResolutionMode"`
	DeviceName        string  `json:"deviceName"`
	Bright            float32 `json:"bright"`
	MusicVolume       float32 `json:"musicVolume"`
	EffectVolume      float32 `json:"effectVolume"`
	HitSoundVolume    float32 `json:"hitSoundVolume"`
	SoundOffset       float32 `json:"soundOffset"`
	NoteScale         float32 `json:"noteScale"`
}

type User struct {
	ShowPlayerId bool   `json:"showPlayerId"`
	SelfIntro    string `json:"selfIntro"`
	Avatar       string `json:"avatar"`
	Background   string `json:"background"`
}

type ScoreAcc struct {
	Score      int     `json:"score"`
	Acc        float32 `json:"acc"`
	Level      string  `json:"level"`
	Fc         bool    `json:"fc"`
	SongId     string  `json:"songId"`
	Difficulty float32 `json:"difficulty"`
	Rks        float32 `json:"rks"`
}

// net struct
type UserMe struct {
	ACL struct {
		NAMING_FAILED struct {
			Write bool `json:"write"`
			Read  bool `json:"read"`
		} `json:"*"`
	} `json:"ACL"`
	AuthData struct {
		Taptap struct {
			AccessToken  string `json:"access_token"`
			Avatar       string `json:"avatar"`
			Kid          string `json:"kid"`
			MacAlgorithm string `json:"mac_algorithm"`
			MacKey       string `json:"mac_key"`
			Name         string `json:"name"`
			Openid       string `json:"openid"`
			TokenType    string `json:"token_type"`
			Unionid      string `json:"unionid"`
		} `json:"taptap"`
	} `json:"authData"`
	Avatar              string    `json:"avatar"`
	CreatedAt           time.Time `json:"createdAt"`
	EmailVerified       bool      `json:"emailVerified"`
	MobilePhoneVerified bool      `json:"mobilePhoneVerified"`
	Nickname            string    `json:"nickname"`
	ObjectID            string    `json:"objectId"`
	SessionToken        string    `json:"sessionToken"`
	ShortID             string    `json:"shortId"`
	UpdatedAt           time.Time `json:"updatedAt"`
	Username            string    `json:"username"`
}

type GameSave struct {
	Results []struct {
		CreatedAt time.Time `json:"createdAt"`
		GameFile  struct {
			Type      string    `json:"__type"`
			Bucket    string    `json:"bucket"`
			CreatedAt time.Time `json:"createdAt"`
			Key       string    `json:"key"`
			MetaData  struct {
				Checksum string `json:"_checksum"`
				Prefix   string `json:"prefix"`
				Size     int    `json:"size"`
			} `json:"metaData"`
			MimeType  string    `json:"mime_type"`
			Name      string    `json:"name"`
			ObjectID  string    `json:"objectId"`
			Provider  string    `json:"provider"`
			UpdatedAt time.Time `json:"updatedAt"`
			URL       string    `json:"url"`
		} `json:"gameFile"`
		ModifiedAt struct {
			Type string    `json:"__type"`
			Iso  time.Time `json:"iso"`
		} `json:"modifiedAt"`
		Name      string    `json:"name"`
		ObjectID  string    `json:"objectId"`
		Summary   string    `json:"summary"`
		UpdatedAt time.Time `json:"updatedAt"`
		User      struct {
			Type      string `json:"__type"`
			ClassName string `json:"className"`
			ObjectID  string `json:"objectId"`
		} `json:"user"`
	} `json:"results"`
}
type PlayerInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Avatar    string    `json:"avatar"`
}
type UserRecord struct {
	Session    string      `json:"session"`
	PlayerInfo *PlayerInfo `json:"playerInfo"`
	ScoreAcc   []ScoreAcc  `json:"scoreAcc"`
	Summary    *Summary    `json:"summary"`
}

type RespCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Summary struct {
	SaveVersion       byte
	ChallengeModeRank int16
	Rks               float32
	GameVersion       byte
	Avatar            string
	ScoreAcc          [4]SummaryScoreAcc
}
type SummaryScoreAcc struct {
	Cleared   int16
	FullCombo int16
	Phi       int16
}
