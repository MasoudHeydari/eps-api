package delivery

import (
	"context"
	"time"

	"github.com/MasoudHeydari/eps-api/db"
	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/MasoudHeydari/eps-api/model"
	"github.com/sirupsen/logrus"
)

func (s *Server) PollJob() {
	ctx := context.Background()
	for {
		time.Sleep(time.Minute * 2)
		jobID, sqID, err := db.GetAJobID(ctx, s.db)
		if err != nil {
			switch {
			case ent.IsNotFound(err):
				logrus.Info(
					"server.PollJob.GetJobID.IsNotFound: there is no active sq-id to fetch its data",
					"error", err,
				)
			default:
				logrus.Info("server.PollJob.GetJobID: error getting job id", err)
			}
			continue
		}
		items, shouldCancelTheJob, err := s.agent.PollJob(jobID)
		if err != nil {
			logrus.Info("server.PollJob.agent.PollJob: error fetching api responses", err)
			if shouldCancelTheJob {
				err = db.CancelSQ(ctx, s.db, sqID)
				if err != nil {
					logrus.Info(
						"server.PollJob.CancelSQ: failed to finish sqid",
						"sq_id", sqID,
						"error", err,
					)
				}
			}
			continue
		}
		for _, item := range items {
			// extract emails
			linkURL := item.URL
			emails, err := s.agent.extractEmails(linkURL)
			if err != nil {
				logrus.Errorf("Search: %v", err)
			}

			// extract phones
			phones, err := s.agent.extractPhoneNumbersFromAllPossibleURLs(linkURL)
			if err != nil {
				logrus.Errorf("Search: %v", err)
			}

			// extract key-words
			var keyWords []string
			keyWords, err = s.agent.extractKeywords(linkURL)
			if err != nil {
				logrus.Errorf("Search: %v", err)
			}

			result := model.SearchResult{
				Rank:        item.RankAbsolute,
				URL:         item.URL,
				Title:       item.Title,
				Phones:      phones,
				Emails:      emails,
				KeyWords:    keyWords,
				Description: item.Description,
			}
			err = db.InsertSearchResult(ctx, s.db, result, sqID)
			if err != nil {
				logrus.Info("server.PollJob.InsertSearchResult: ", err)
			}
		}
		err = db.CancelSQ(ctx, s.db, sqID)
		if err != nil {
			logrus.Info(
				"server.PollJob.CancelSQ: failed to finish sqid",
				"sq_id", sqID,
				"error", err,
			)
		}
	}
}
