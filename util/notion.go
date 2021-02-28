package util

import (
	"log"

	"github.com/kjk/notionapi"
)

func NotionClient(authToken string) *notionapi.Client {
	client := &notionapi.Client{}
	client.AuthToken = authToken
	return client
}

func getPages(client *notionapi.Client, mainPageID string) []string {
	page, err := client.DownloadPage(mainPageID)
	if err != nil {
		panic(err)
	}
	return page.GetSubPages()
}

//Download page and get the page information
/*
type Page struct {
	ID string

	// expose raw records for all data associated with this page
	BlockRecords          []*Record
	UserRecords           []*Record
	CollectionRecords     []*Record
	CollectionViewRecords []*Record
	DiscussionRecords     []*Record
	CommentRecords        []*Record

	// for every block of type collection_view and its view_ids
	// we } TableView representing that collection view_id
	TableViews []*TableView
	// contains filtered or unexported fields
}
*/
func getPage(client *notionapi.Client, pageID string) notionapi.Page {
	page, err := client.DownloadPage(pageID)
	if err != nil {
		panic(err)
	}
	return *page
}

func NotionPages(client *notionapi.Client, mainPageID string) {
	// list subpage in Posts database page
	pages := getPages(client, mainPageID)
	for _, pageID := range pages {
		fm := frontMatter{PageID: pageID}
		page := getPage(client, pageID)
		page.ForEachBlock(func(block *notionapi.Block) {
			fm.Check(block)
		})

		mdfile := fm.GetFile()
		log.Println(mdfile)
		fmtMediaLink(mdfile)
		fmtFrontMatter(mdfile, fm)
	}
}
