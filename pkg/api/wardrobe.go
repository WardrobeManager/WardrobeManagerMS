//
// definitions.go
//
// May 2021, Prashant Desai
//

package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type WardrobeService interface {
	AddWardrobe(new NewWardrobeRequest) error
	DeleteWardrobe(user string, id string) error
	GetWardrobe(user string, id string) (*GetWardrobeResponse, error)
	GetAllWardrobe(user string) ([]*GetWardrobeResponse, error)
	GetFile(filename string, cbHandler HandleFile) error

	AddOutfit(new NewOutfitRequest) error
	DeleteOutfit(user string, id string) error
	GetOutfit(user string, id string) (*GetOutfitResponse, error)
	GetAllOutfits(user string) ([]*GetOutfitResponse, error)
}

type WardrobeRepository interface {
	Add(user string, wardrobes *WardrobeCloset) error
	Get(user string) (*WardrobeCloset, error)
	Update(user string, wardrobes *WardrobeCloset) error
	DeleteAll(user string) error
}

type ImageRepository interface {
	AddFile(name string, file []byte) error
	GetFile(name string) ([]byte, error)
	UpdateFile(name string, file []byte) error
	DeleteFile(name string) error
	AddFileFromFile(name string, rd io.Reader) error
	GetFileWithHandler(filename string, fileHandler HandleFile) error
}

type wardrobeService struct {
	mu      sync.Mutex
	db      WardrobeRepository
	imageDb ImageRepository
	l       *wardrobeLabelToText
}

func NewWardrobeService(dbIn WardrobeRepository, imageDbIn ImageRepository, rds, rx, tx string) (WardrobeService, error) {

	glog.Infof("Creating Wardrobe Service")

	service := &wardrobeService{
		db:      dbIn,
		imageDb: imageDbIn,
	}

	var err error
	service.l, err = newWardrobeLabelToText(rds, rx, tx, service)
	if err != nil {
		glog.Errorf("error initializing label to text service endpoint : {err=%v}", err)
		return nil, err
	}

	go service.l.receiveLoop()

	return service, nil
}

func (w *wardrobeService) AddWardrobe(newWd NewWardrobeRequest) error {

	var addUser bool = false

	// generate a unique id
	id := uuid.New().String()

	glog.Infof("adding wardrobe {user=%s}, {id=%s}", newWd.User, id)

	w.mu.Lock()
	defer w.mu.Unlock()

	wc, err := w.db.Get(newWd.User)
	switch err := err.(type) {
	case nil:
	case *UserNotFound:
		addUser = true
		wc = &WardrobeCloset{
			User:      newWd.User,
			Wardrobes: make([]Wardrobe, 0),
			Outfits:   make([]Outfit, 0),
		}
	case *ResourceUnavailable:
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return fmt.Errorf("Unknown error : %w", err)
	}

	//Store image to file
	imageFile := genUniqImageFileName(newWd.User, id)
	labelFile := genUniqLabelFileName(newWd.User, id)

	for _, file := range []string{imageFile, labelFile} {
		_, err := w.imageDb.GetFile(file)
		switch err := err.(type) {
		case NoSuchFileOrDirectory:
			// this is good
			break
		case nil:
			// file with same name found
			return &DuplicateFile{
				File: file,
			}
		default:
			return fmt.Errorf("File system access error : %w", err)
		}
	}

	//Store files
	mimeMainFile, err1 := newWd.MainImageMime.Open()
	if err1 != nil {
		return fmt.Errorf("Error opening mime main image file : %w", err1)
	}
	defer mimeMainFile.Close()

	mimeLabelFile, err2 := newWd.LabelImageMime.Open()
	if err2 != nil {
		return fmt.Errorf("Error opening mime label image file : %w", err2)
	}
	defer mimeLabelFile.Close()

	err = w.imageDb.AddFileFromFile(imageFile, mimeMainFile)
	if err != nil {
		return fmt.Errorf("Error saving image to file system : %w", err)
	}

	err = w.imageDb.AddFileFromFile(labelFile, mimeLabelFile)
	if err != nil {
		return fmt.Errorf("Error saving image to file system : %w", err)
	}

	//Update user
	wc.Wardrobes = append(wc.Wardrobes, Wardrobe{
		Identifier:  id,
		MainFile:    imageFile,
		LabelFile:   labelFile,
		Description: newWd.Description,
	})
	if addUser == true {
		err = w.db.Add(newWd.User, wc)
	} else {
		err = w.db.Update(newWd.User, wc)
	}
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	//label to text
	sEnc := base64.StdEncoding.EncodeToString(newWd.LabelImage)
	err = w.l.sendLabel(newWd.User, id, sEnc)
	if err != nil {
		glog.Warningf("failure while trying to send lable from label to text {err=%v}", err)
	}

	glog.Infof("done adding wardrobe {user=%s}, {id=%s}", newWd.User, id)

	return nil
}

func (w *wardrobeService) DeleteWardrobe(user string, id string) error {

	glog.Infof("deleting wardrobe {user=%s}, {id=%s}", user, id)

	w.mu.Lock()
	defer w.mu.Unlock()

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		return fmt.Errorf("Empty closet")
	}

	tmp := wc.Wardrobes[:0]
	for _, ward := range wc.Wardrobes {
		if ward.Identifier == id {
			err = w.imageDb.DeleteFile(ward.MainFile)
			if err != nil {
				fmt.Printf("Error deleting image file : %w", err)
			}

			err = w.imageDb.DeleteFile(ward.LabelFile)
			if err != nil {
				fmt.Printf("Error deleting label file : %w", err)
			}
		} else {
			tmp = append(tmp, ward)
		}
	}
	wc.Wardrobes = tmp

	err = w.db.Update(user, wc)
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	glog.Infof("done deleting wardrobe {user=%s}, {id=%s}", user, id)

	return nil
}

func (w *wardrobeService) GetWardrobe(user string, id string) (*GetWardrobeResponse, error) {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return nil, fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return nil, fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return nil, fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		return nil, fmt.Errorf("Empty closet")
	}

	for _, ward := range wc.Wardrobes {
		if ward.Identifier == id {
			wardReq := &GetWardrobeResponse{
				Id:          ward.Identifier,
				Description: ward.Description,
				MainImage:   ward.MainFile,
				LabelImage:  ward.LabelFile,
			}

			return wardReq, nil
		}
	}

	return nil, UserNotFound{User: user}
}

func (w *wardrobeService) GetAllWardrobe(user string) ([]*GetWardrobeResponse, error) {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return nil, fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return nil, fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return nil, fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		return nil, fmt.Errorf("Empty closet")
	}

	wardReqs := make([]*GetWardrobeResponse, 0)
	for _, ward := range wc.Wardrobes {
		wardReq := &GetWardrobeResponse{
			Id:          ward.Identifier,
			Description: ward.Description,
			MainImage:   ward.MainFile,
			LabelImage:  ward.LabelFile,
		}

		wardReqs = append(wardReqs, wardReq)
	}

	return wardReqs, nil
}

func (w *wardrobeService) GetFile(filename string, cb HandleFile) error {

	return w.imageDb.GetFileWithHandler(filename, cb)
}

func (w *wardrobeService) AddOutfit(newOt NewOutfitRequest) error {

	// generate a unique id
	id := uuid.New().String()

	glog.Infof("adding outfit {user=%s}, {id=%s}", newOt.User, id)

	w.mu.Lock()
	defer w.mu.Unlock()

	wc, err := w.db.Get(newOt.User)
	switch err := err.(type) {
	case nil:
	case *UserNotFound:
		return fmt.Errorf("User not found  : %w", err)
	case *ResourceUnavailable:
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return fmt.Errorf("Unknown error : %w", err)
	}

	//Update user
	wc.Outfits = append(wc.Outfits, Outfit{
		Identifier:   id,
		TopId:        newOt.TopId,
		BottomId:     newOt.BottomId,
		Description:  newOt.Description,
		LikeCount:    0,
		DislikeCount: 0,
	})

	err = w.db.Update(newOt.User, wc)
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	glog.Infof("done adding outfit {user=%s}, {id=%s}", newOt.User, id)

	return nil
}

func (w *wardrobeService) DeleteOutfit(user string, id string) error {

	glog.Infof("deleting outfit {user=%s}, {id=%s}", user, id)

	w.mu.Lock()
	defer w.mu.Unlock()

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Outfits) == 0 {
		return fmt.Errorf("Empty outfits")
	}

	tmp := wc.Outfits[:0]
	for _, ot := range wc.Outfits {
		if ot.Identifier != id {
			tmp = append(tmp, ot)
		}
	}
	wc.Outfits = tmp

	err = w.db.Update(user, wc)
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	glog.Infof("done deleting outfit {user=%s}, {id=%s}", user, id)

	return nil
}

func (w *wardrobeService) GetOutfit(user string, id string) (*GetOutfitResponse, error) {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return nil, fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return nil, fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return nil, fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Outfits) == 0 {
		return nil, fmt.Errorf("Empty closet")
	}

	for _, ot := range wc.Outfits {
		if ot.Identifier == id {
			otReq := &GetOutfitResponse{
				Id:           ot.Identifier,
				TopId:        ot.TopId,
				BottomId:     ot.BottomId,
				Description:  ot.Description,
				LikeCount:    ot.LikeCount,
				DislikeCount: ot.DislikeCount,
			}

			return otReq, nil
		}
	}

	return nil, UserNotFound{User: user}
}

func (w *wardrobeService) GetAllOutfits(user string) ([]*GetOutfitResponse, error) {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return nil, fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return nil, fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return nil, fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Outfits) == 0 {
		return nil, fmt.Errorf("No outfits")
	}

	otReqs := make([]*GetOutfitResponse, 0)
	for _, ot := range wc.Outfits {
		otReq := &GetOutfitResponse{
			Id:           ot.Identifier,
			TopId:        ot.TopId,
			BottomId:     ot.BottomId,
			Description:  ot.Description,
			LikeCount:    ot.LikeCount,
			DislikeCount: ot.DislikeCount,
		}

		otReqs = append(otReqs, otReq)
	}

	return otReqs, nil
}

//private functions
func (w *wardrobeService) updateWardrobeLabelText(user, id, text string) error {

	glog.Infof("Received text {user=%s}, {id=%s}, {text=%s}", user, id, text)

	w.mu.Lock()
	defer w.mu.Unlock()

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		glog.Errorf("User not found {user=%s}, {err=%v}", user, err)
		return fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		glog.Errorf("Wardrobe db is unavailable : %w", err)
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		glog.Errorf("Unknown error : %w", err)
		return fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		glog.Errorf("Empty closet")
		return fmt.Errorf("Empty closet")
	}

	tmp := wc.Wardrobes[:0]
	for _, ward := range wc.Wardrobes {
		if ward.Identifier == id {
			ward.LabelText = text
		}

		tmp = append(tmp, ward)
	}
	wc.Wardrobes = tmp

	err = w.db.Update(user, wc)
	switch err := err.(type) {
	case nil:
	default:
		glog.Errorf("Database access failure : %w", err)
		return fmt.Errorf("Database access failure : %w", err)
	}

	return nil
}

//label to text service
type wardrobeLabelToText struct {
	conn      redis.Conn
	psc       redis.PubSubConn
	txconn    redis.Conn
	rxChannel string
	txChannel string
	s         *wardrobeService
}

func newWardrobeLabelToText(redisServerAddr, rxChannel, txChannel string, s *wardrobeService) (*wardrobeLabelToText, error) {

	glog.Infof("Dialing redis server : {address=%s}", redisServerAddr)

	c, err := redis.Dial("tcp", redisServerAddr)
	if err != nil {
		glog.Errorf("error dialing redis server to establish RX conn: {err=%v}", err)
		return nil, err
	}

	txc, err := redis.Dial("tcp", redisServerAddr)
	if err != nil {
		glog.Errorf("error dialing redis server to establish TX conn: {err=%v}", err)
		return nil, err
	}

	l := &wardrobeLabelToText{
		conn: c,
		psc: redis.PubSubConn{
			Conn: c,
		},
		txconn:    txc,
		rxChannel: rxChannel,
		txChannel: txChannel,
		s:         s,
	}

	return l, nil

}

func (s *wardrobeLabelToText) receiveLoop() {
	glog.Infof("Running label to text receive loop")

	if err := s.psc.Subscribe(redis.Args{}.AddFlat(s.rxChannel)...); err != nil {
		glog.Errorf("Error subscribing to receive {channel=%s} : {err=%v}", s.rxChannel, err)
		return
	}

	for {
		switch n := s.psc.Receive().(type) {
		case error:
			glog.Errorf("Received error from redis server {err=%v}", n)
			return
		case redis.Message:
			if err := s.onMessageReceive(n.Channel, n.Data); err != nil {
				glog.Errorf("Error processing message received on {channel=%s}", n.Channel)
			}
		case redis.Subscription:
			switch n.Count {
			case 1:
				glog.Infof("Subscribe to {channel=%s}", n.Channel)
			case 0:
				glog.Errorf("Unexpected unsubscribe to {channel=%s}", n.Channel)
				return
			}
		}
	}
}

func (s *wardrobeLabelToText) onMessageReceive(channel string, data []byte) error {
	glog.Infof("Received message on {channel=%s}, {size=%d}", channel, len(data))

	var resp LabelToTextResponse
	err1 := json.Unmarshal(data, &resp)
	if err1 != nil {
		glog.Error("error unmarshaling received label to text json response {err=%v}", err1)
		return err1
	}

	err2 := s.s.updateWardrobeLabelText(resp.User, resp.Id, resp.Text)
	if err2 != nil {
		glog.Error("error updating wardrobe label text  {err=%v}", err2)
		return err2
	}

	return nil
}

func (s *wardrobeLabelToText) sendLabel(user, id, image string) error {

	req := &LabelToTextRequest{
		User:     user,
		Id:       id,
		RawImage: image,
	}

	jsonReq, err3 := json.Marshal(req)
	if err3 != nil {
		glog.Error("error marshaling text json output {err=%v}", err3)
		return err3
	}

	if _, err := s.txconn.Do("PUBLISH", s.txChannel, jsonReq); err != nil {
		glog.Errorf("error publishing label to text request to tx channel, {txChannel=%s}, {err=%v}", s.txChannel, err)
		return err
	}
	return nil
}

// Error codes
func (e UserNotFound) Error() string {
	return fmt.Sprintf("User %s not found", e.User)
}

func (e NoSuchFileOrDirectory) Error() string {
	return fmt.Sprintf("File %s not found", e.File)
}

func (e ResourceUnavailable) Error() string {
	return fmt.Sprintf("Service %s is down", e.Server)
}

func (e DuplicateFile) Error() string {
	return fmt.Sprintf("Duplicate file name %s", e.File)
}

/*
func (e DuplicateFile) Is(target error) bool {
	switch target.(type) {
	default:
		return false
	case *DuplicateFile:
		return true
	}

}
*/

func genUniqImageFileName(user string, filename string) string {
	stringToHash := []byte(user + "_image_" + filename)
	md5Bytes := md5.Sum(stringToHash)
	return hex.EncodeToString(md5Bytes[:])
}

func genUniqLabelFileName(user string, filename string) string {
	stringToHash := []byte(user + "_label_" + filename)
	md5Bytes := md5.Sum(stringToHash)
	return hex.EncodeToString(md5Bytes[:])

}
