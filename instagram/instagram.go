package instagram

import (
	"github.com/patrickmn/go-cache"
	"github.com/siongui/instago"
	log "github.com/sirupsen/logrus"
	"instabot/config"
	"strconv"
	"time"
)

type InstaManager struct {
	privateAPIManager instago.IGApiManager
	publicAPIManager  instago.IGApiManager
	infoCache         *cache.Cache
	storiesCache      *cache.Cache
}

type AccountInfo struct {
	UserId    string
	Username  string
	Biography string
	Type      string
	Followers string
}

const (
	PrivateType = "private"
	PublicType  = "public"
)

func NewInstagramManager() InstaManager {
	privateAPIManager, err := instago.NewInstagramApiManager(config.Cfg.Instagram.AuthLocation)
	if err != nil {
		panic("instagram manager was not loaded")
	}
	return InstaManager{
		privateAPIManager: *privateAPIManager,
		publicAPIManager:  *instago.NewApiManager(nil, nil),
		infoCache:         cache.New(5*time.Minute, 10*time.Minute),
		storiesCache:      cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (manager InstaManager) GetUserInfo(username string) *AccountInfo {
	if cachedInfo, found := manager.infoCache.Get(username); found {
		return cachedInfo.(*AccountInfo)
	}
	info, err := manager.privateAPIManager.GetUserInfo(username)
	if err != nil {
		log.Errorf("Can't get info for username %s: %v", username, err)
		return nil
	}
	accountType := PrivateType
	if info.IsPublic() {
		accountType = PublicType
	}
	accountInfo := &AccountInfo{
		UserId:    info.Id,
		Username:  info.Username,
		Biography: info.Biography,
		Type:      accountType,
		Followers: strconv.FormatInt(info.EdgeFollowedBy.Count, 10),
	}
	_ = manager.infoCache.Add(username, accountInfo, cache.DefaultExpiration)
	return accountInfo
}

func (manager InstaManager) GetRecentPosts(username string) ([]instago.IGMedia, error) {
	return manager.privateAPIManager.GetRecentPostMedia(username)
}

func (manager InstaManager) GetStories(username string) []instago.IGItem {
	if cachedStories, found := manager.storiesCache.Get(username); found {
		return cachedStories.([]instago.IGItem)
	}
	info := manager.GetUserInfo(username)
	if info == nil {
		log.Errorf("Can't get stories for username %s", username)
		return make([]instago.IGItem, 0)
	}
	story, err := manager.privateAPIManager.GetUserStory(info.UserId)
	if err != nil {
		log.Errorf("Can't get stories for username %s: %v", username, err)
		return make([]instago.IGItem, 0)
	}
	_ = manager.storiesCache.Add(username, story.Reel.GetItems(), cache.DefaultExpiration)
	return story.Reel.GetItems()
}
