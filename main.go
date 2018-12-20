package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	//"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"sync"

	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
	"net"
	"io"
)

//"show_suggested_profiles":false,"graphql":
//,"edge_saved_media":{"count":0,"page_info":{"has_next_page":false,"end_cursor":null},"edges":[]},"edge_media_collections":{"count":0,"page_info":{"has_next_page":false,"end_cursor":null},"edges":[]}}}
type sharedData struct {
	ActivityCounts struct {
		CommentLikes  int `json:"comment_likes"`
		Comments      int `json:"comments"`
		Likes         int `json:"likes"`
		Relationships int `json:"relationships"`
		Usertags      int `json:"usertags"`
	} `json:"activity_counts"`
	Config struct {
		CsrfToken string `json:"csrf_token"`
		Viewer    struct {
			AllowContactsSync bool   `json:"allow_contacts_sync"`
			Biography         string `json:"biography"`
			ExternalURL       string `json:"external_url"`
			FullName          string `json:"full_name"`
			HasProfilePic     bool   `json:"has_profile_pic"`
			ID                string `json:"id"`
			ProfilePicURL     string `json:"profile_pic_url"`
			ProfilePicURLHd   string `json:"profile_pic_url_hd"`
			Username          string `json:"username"`
		} `json:"viewer"`
	} `json:"config"`
	SupportsEs6  bool   `json:"supports_es6"`
	CountryCode  string `json:"country_code"`
	LanguageCode string `json:"language_code"`
	Locale       string `json:"locale"`
	EntryData    struct {
		ProfilePage []struct {
			LoggingPageID         string `json:"logging_page_id"`
			ShowSuggestedProfiles bool   `json:"show_suggested_profiles"`
			Graphql               struct {
				User struct {
					Biography              string      `json:"biography"`
					BlockedByViewer        bool        `json:"blocked_by_viewer"`
					CountryBlock           bool        `json:"country_block"`
					ExternalURL            interface{} `json:"external_url"`
					ExternalURLLinkshimmed interface{} `json:"external_url_linkshimmed"`
					EdgeFollowedBy         struct {
						Count int `json:"count"`
					} `json:"edge_followed_by"`
					FollowedByViewer bool `json:"followed_by_viewer"`
					EdgeFollow       struct {
						Count int `json:"count"`
					} `json:"edge_follow"`
					FollowsViewer      bool   `json:"follows_viewer"`
					FullName           string `json:"full_name"`
					HasChannel         bool   `json:"has_channel"`
					HasBlockedViewer   bool   `json:"has_blocked_viewer"`
					HighlightReelCount int    `json:"highlight_reel_count"`
					HasRequestedViewer bool   `json:"has_requested_viewer"`
					ID                 string `json:"id"`
					IsPrivate          bool   `json:"is_private"`
					IsVerified         bool   `json:"is_verified"`
					MutualFollowers    struct {
						AdditionalCount int      `json:"additional_count"`
						Usernames       []string `json:"usernames"`
					} `json:"mutual_followers"`
					ProfilePicURL                string      `json:"profile_pic_url"`
					ProfilePicURLHd              string      `json:"profile_pic_url_hd"`
					RequestedByViewer            bool        `json:"requested_by_viewer"`
					Username                     string      `json:"username"`
					ConnectedFbPage              interface{} `json:"connected_fb_page"`
					EdgeFelixCombinedPostUploads struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_felix_combined_post_uploads"`
					EdgeFelixCombinedDraftUploads struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_felix_combined_draft_uploads"`
					EdgeFelixVideoTimeline struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_felix_video_timeline"`
					EdgeFelixDrafts struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_felix_drafts"`
					EdgeFelixPendingPostUploads struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_felix_pending_post_uploads"`
					EdgeFelixPendingDraftUploads struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_felix_pending_draft_uploads"`
					EdgeOwnerToTimelineMedia struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool   `json:"has_next_page"`
							EndCursor   string `json:"end_cursor"`
						} `json:"page_info"`
						Edges []struct {
							Node struct {
								Typename           string `json:"__typename"`
								ID                 string `json:"id"`
								EdgeMediaToCaption struct {
									Edges []struct {
										Node struct {
											Text string `json:"text"`
										} `json:"node"`
									} `json:"edges"`
								} `json:"edge_media_to_caption"`
								Shortcode          string `json:"shortcode"`
								EdgeMediaToComment struct {
									Count int `json:"count"`
								} `json:"edge_media_to_comment"`
								CommentsDisabled bool `json:"comments_disabled"`
								TakenAtTimestamp int  `json:"taken_at_timestamp"`
								Dimensions       struct {
									Height int `json:"height"`
									Width  int `json:"width"`
								} `json:"dimensions"`
								DisplayURL  string `json:"display_url"`
								EdgeLikedBy struct {
									Count int `json:"count"`
								} `json:"edge_liked_by"`
								EdgeMediaPreviewLike struct {
									Count int `json:"count"`
								} `json:"edge_media_preview_like"`
								GatingInfo   interface{} `json:"gating_info"`
								MediaPreview string      `json:"media_preview"`
								Owner        struct {
									ID string `json:"id"`
								} `json:"owner"`
								ThumbnailSrc       string `json:"thumbnail_src"`
								ThumbnailResources []struct {
									Src          string `json:"src"`
									ConfigWidth  int    `json:"config_width"`
									ConfigHeight int    `json:"config_height"`
								} `json:"thumbnail_resources"`
								IsVideo        bool `json:"is_video"`
								VideoViewCount int  `json:"video_view_count"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_owner_to_timeline_media"`
					EdgeSavedMedia struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_saved_media"`
					EdgeMediaCollections struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool        `json:"has_next_page"`
							EndCursor   interface{} `json:"end_cursor"`
						} `json:"page_info"`
						Edges []interface{} `json:"edges"`
					} `json:"edge_media_collections"`
				} `json:"user"`
			} `json:"graphql"`
			FelixOnboardingVideoResources struct {
				Mp4    string `json:"mp4"`
				Poster string `json:"poster"`
			} `json:"felix_onboarding_video_resources"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
	Gatekeepers struct {
		Cb     bool `json:"cb"`
		Ld     bool `json:"ld"`
		Vl     bool `json:"vl"`
		Seo    bool `json:"seo"`
		Seoht  bool `json:"seoht"`
		TwoFac bool `json:"2fac"`
		Sf     bool `json:"sf"`
		Saa    bool `json:"saa"`
	} `json:"gatekeepers"`
	Knobs struct {
		AcctNtb int `json:"acct:ntb"`
		Cb      int `json:"cb"`
		Captcha int `json:"captcha"`
	} `json:"knobs"`
	Qe struct {
		FormNavigationDialog struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"form_navigation_dialog"`
		CredMan struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"cred_man"`
		DashForVod struct {
			G string `json:"g"`
			P struct {
				EnforceGl            string `json:"enforce_gl"`
				IsEnabled            string `json:"is_enabled"`
				IsFeedOrganicEnabled string `json:"is_feed_organic_enabled"`
				Variant              string `json:"variant"`
			} `json:"p"`
		} `json:"dash_for_vod"`
		ProfileHeaderName struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"profile_header_name"`
		Bc3L struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"bc3l"`
		DirectConversationReporting struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"direct_conversation_reporting"`
		GeneralReporting struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"general_reporting"`
		Reporting struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"reporting"`
		AccRecoveryLink struct {
			G string `json:"g"`
			P struct {
				ShowAccountRecoveryRedesign   string `json:"show_account_recovery_redesign"`
				ShowResetPasswordInterstitial string `json:"show_reset_password_interstitial"`
			} `json:"p"`
		} `json:"acc_recovery_link"`
		Notif struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"notif"`
		FbUnlink struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"fb_unlink"`
		MobileStoriesDoodling struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"mobile_stories_doodling"`
		MoveCommentInputToTop struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"move_comment_input_to_top"`
		MobileCancel struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"mobile_cancel"`
		MobileSearchRedesign struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"mobile_search_redesign"`
		ShowCopyLink struct {
			G string `json:"g"`
			P struct {
				ShowCopyLinkOption string `json:"show_copy_link_option"`
			} `json:"p"`
		} `json:"show_copy_link"`
		MobileLogout struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"mobile_logout"`
		PEdit struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"p_edit"`
		Four04AsReact struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"404_as_react"`
		AccRecovery struct {
			G string `json:"g"`
			P struct {
				HasPrefill string `json:"has_prefill"`
			} `json:"p"`
		} `json:"acc_recovery"`
		Collections struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"collections"`
		CommentTa struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"comment_ta"`
		Connections struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"connections"`
		DiscPpl struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"disc_ppl"`
		EbdUl struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"ebd_ul"`
		EbdsimLi struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"ebdsim_li"`
		EbdsimLo struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"ebdsim_lo"`
		EmptyFeed struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"empty_feed"`
		Bundles struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"bundles"`
		ExitStoryCreation struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"exit_story_creation"`
		GdprLoggedOut struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"gdpr_logged_out"`
		Appsell struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"appsell"`
		Imgopt struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"imgopt"`
		FollowButton struct {
			G string `json:"g"`
			P struct {
				IsInline string `json:"is_inline"`
			} `json:"p"`
		} `json:"follow_button"`
		Loggedout struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"loggedout"`
		LoggedoutUpsell struct {
			G string `json:"g"`
			P struct {
				HasNewLoggedoutUpsellContent string `json:"has_new_loggedout_upsell_content"`
			} `json:"p"`
		} `json:"loggedout_upsell"`
		Msisdn struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"msisdn"`
		BgSync struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"bg_sync"`
		Onetaplogin struct {
			G string `json:"g"`
			P struct {
				AfterLogin     string `json:"after_login"`
				StorageVersion string `json:"storage_version"`
			} `json:"p"`
		} `json:"onetaplogin"`
		LoginPoe struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"login_poe"`
		PrivateLo struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"private_lo"`
		ProfileTabs struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"profile_tabs"`
		PushNotifications struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"push_notifications"`
		Reg struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"reg"`
		RegVp struct {
			G string `json:"g"`
			P struct {
				HideValueProp string `json:"hide_value_prop"`
			} `json:"p"`
		} `json:"reg_vp"`
		ReportMedia struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"report_media"`
		ReportProfile struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"report_profile"`
		SidecarSwipe struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"sidecar_swipe"`
		SuUniverse struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"su_universe"`
		Stale struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"stale"`
		StoriesLo struct {
			G string `json:"g"`
			P struct {
				ContextualLogin string `json:"contextual_login"`
			} `json:"p"`
		} `json:"stories_lo"`
		Stories struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"stories"`
		TpPblshr struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"tp_pblshr"`
		Video struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"video"`
		GdprSettings struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"gdpr_settings"`
		GdprEuTos struct {
			G string `json:"g"`
			P struct {
				GdprRequired  string `json:"gdpr_required"`
				EuNewUserFlow string `json:"eu_new_user_flow"`
				TosVersion    string `json:"tos_version"`
			} `json:"p"`
		} `json:"gdpr_eu_tos"`
		GdprRowTos struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"gdpr_row_tos"`
		FdGr struct {
			G string `json:"g"`
			P struct {
				ShowFollowToast string `json:"show_follow_toast"`
			} `json:"p"`
		} `json:"fd_gr"`
		Felix struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"felix"`
		FelixClearFbCookie struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
				Blacklist string `json:"blacklist"`
			} `json:"p"`
		} `json:"felix_clear_fb_cookie"`
		FelixCreationDurationLimits struct {
			G string `json:"g"`
			P struct {
				MinimumLengthSeconds string `json:"minimum_length_seconds"`
				MaximumLengthSeconds string `json:"maximum_length_seconds"`
			} `json:"p"`
		} `json:"felix_creation_duration_limits"`
		FelixCreationEnabled struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"felix_creation_enabled"`
		FelixCreationFbCrossposting struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"felix_creation_fb_crossposting"`
		FelixCreationFbCrosspostingV2 struct {
			G string `json:"g"`
			P struct {
				IsEnabled string `json:"is_enabled"`
			} `json:"p"`
		} `json:"felix_creation_fb_crossposting_v2"`
		FelixCreationValidation struct {
			G string `json:"g"`
			P struct {
				EditVideoControls        string `json:"edit_video_controls"`
				MaxVideoSizeInBytes      string `json:"max_video_size_in_bytes"`
				TitleMaximumLength       string `json:"title_maximum_length"`
				DescriptionMaximumLength string `json:"description_maximum_length"`
				ValidCoverMimeTypes      string `json:"valid_cover_mime_types"`
				ValidVideoMimeTypes      string `json:"valid_video_mime_types"`
				ValidVideoExtensions     string `json:"valid_video_extensions"`
			} `json:"p"`
		} `json:"felix_creation_validation"`
		FelixCreationVideoUpload struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"felix_creation_video_upload"`
		FelixEarlyOnboarding struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"felix_early_onboarding"`
		Pride struct {
			G string `json:"g"`
			P struct {
				Enabled          string `json:"enabled"`
				HashtagWhitelist string `json:"hashtag_whitelist"`
			} `json:"p"`
		} `json:"pride"`
		UnfollowConfirm struct {
			G string `json:"g"`
			P struct {
				NoUnfollowConfirmation string `json:"no_unfollow_confirmation"`
			} `json:"p"`
		} `json:"unfollow_confirm"`
		ProfileEnhanceLi struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"profile_enhance_li"`
		ProfileEnhanceLo struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"profile_enhance_lo"`
		PhoneConfirm struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"phone_confirm"`
		CommentEnhance struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"comment_enhance"`
		MwebMediaChaining struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"mweb_media_chaining"`
		MwebTopicalExplore struct {
			G string `json:"g"`
			P struct {
				ShouldShowQuilt string `json:"should_show_quilt"`
			} `json:"p"`
		} `json:"mweb_topical_explore"`
		WebNametag struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"web_nametag"`
		ImageDowngrade struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"image_downgrade"`
		ImageDowngradeLite struct {
			G string `json:"g"`
			P struct {
				ShouldDowngrade string `json:"should_downgrade"`
			} `json:"p"`
		} `json:"image_downgrade_lite"`
		FollowAllFb struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"follow_all_fb"`
		LiteDirectUpsell struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"lite_direct_upsell"`
		WebLoggedoutNoop struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"web_loggedout_noop"`
		StoriesVideoPreload struct {
			G string `json:"g"`
			P struct {
			} `json:"p"`
		} `json:"stories_video_preload"`
		LiteStoriesVideoPreload struct {
			G string `json:"g"`
			P struct {
				DisablePreload string `json:"disable_preload"`
			} `json:"p"`
		} `json:"lite_stories_video_preload"`
	} `json:"qe"`
	Hostname string `json:"hostname"`
	Platform string `json:"platform"`
	RhxGis   string `json:"rhx_gis"`
	Nonce    string `json:"nonce"`
	ZeroData struct {
	} `json:"zero_data"`
	RolloutHash    string `json:"rollout_hash"`
	BundleVariant  string `json:"bundle_variant"`
	ProbablyHasApp bool   `json:"probably_has_app"`
	ShowAppInstall bool   `json:"show_app_install"`
}

type FollowingJson struct {
	Data struct {
		User struct {
			EdgeFollow struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						ID                string `json:"id"`
						Username          string `json:"username"`
						FullName          string `json:"full_name"`
						ProfilePicURL     string `json:"profile_pic_url"`
						IsVerified        bool   `json:"is_verified"`
						FollowedByViewer  bool   `json:"followed_by_viewer"`
						RequestedByViewer bool   `json:"requested_by_viewer"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_follow"`
		} `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}

type ProfileJson struct {
	LoggingPageID         string `json:"logging_page_id"`
	ShowSuggestedProfiles bool   `json:"show_suggested_profiles"`
	Graphql               struct {
		User struct {
			Biography              string `json:"biography"`
			BlockedByViewer        bool   `json:"blocked_by_viewer"`
			CountryBlock           bool   `json:"country_block"`
			ExternalURL            string `json:"external_url"`
			ExternalURLLinkshimmed string `json:"external_url_linkshimmed"`
			EdgeFollowedBy         struct {
				Count int `json:"count"`
			} `json:"edge_followed_by"`
			FollowedByViewer bool `json:"followed_by_viewer"`
			EdgeFollow       struct {
				Count int `json:"count"`
			} `json:"edge_follow"`
			FollowsViewer      bool   `json:"follows_viewer"`
			FullName           string `json:"full_name"`
			HasChannel         bool   `json:"has_channel"`
			HasBlockedViewer   bool   `json:"has_blocked_viewer"`
			HighlightReelCount int    `json:"highlight_reel_count"`
			HasRequestedViewer bool   `json:"has_requested_viewer"`
			ID                 string `json:"id"`
			IsPrivate          bool   `json:"is_private"`
			IsVerified         bool   `json:"is_verified"`
			MutualFollowers    struct {
				AdditionalCount int           `json:"additional_count"`
				Usernames       []interface{} `json:"usernames"`
			} `json:"mutual_followers"`
			ProfilePicURL                string      `json:"profile_pic_url"`
			ProfilePicURLHd              string      `json:"profile_pic_url_hd"`
			RequestedByViewer            bool        `json:"requested_by_viewer"`
			Username                     string      `json:"username"`
			ConnectedFbPage              interface{} `json:"connected_fb_page"`
			EdgeFelixCombinedPostUploads struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_combined_post_uploads"`
			EdgeFelixCombinedDraftUploads struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_combined_draft_uploads"`
			EdgeFelixVideoTimeline struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_video_timeline"`
			EdgeFelixDrafts struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_drafts"`
			EdgeFelixPendingPostUploads struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_pending_post_uploads"`
			EdgeFelixPendingDraftUploads struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_felix_pending_draft_uploads"`
			EdgeOwnerToTimelineMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int `json:"count"`
						} `json:"edge_media_to_comment"`
						CommentsDisabled bool `json:"comments_disabled"`
						TakenAtTimestamp int  `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int `json:"count"`
						} `json:"edge_media_preview_like"`
						GatingInfo   interface{} `json:"gating_info"`
						MediaPreview string      `json:"media_preview"`
						Owner        struct {
							ID string `json:"id"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo        bool `json:"is_video"`
						VideoViewCount int  `json:"video_view_count"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeSavedMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_saved_media"`
			EdgeMediaCollections struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool        `json:"has_next_page"`
					EndCursor   interface{} `json:"end_cursor"`
				} `json:"page_info"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_collections"`
		} `json:"user"`
	} `json:"graphql"`
	FelixOnboardingVideoResources struct {
		Mp4    string `json:"mp4"`
		Poster string `json:"poster"`
	} `json:"felix_onboarding_video_resources"`
}

type ProfileJsonWindow struct {
	User struct {
		Biography              string `json:"biography"`
		BlockedByViewer        bool   `json:"blocked_by_viewer"`
		CountryBlock           bool   `json:"country_block"`
		ExternalURL            string `json:"external_url"`
		ExternalURLLinkshimmed string `json:"external_url_linkshimmed"`
		EdgeFollowedBy         struct {
			Count int `json:"count"`
		} `json:"edge_followed_by"`
		FollowedByViewer bool `json:"followed_by_viewer"`
		EdgeFollow       struct {
			Count int `json:"count"`
		} `json:"edge_follow"`
		FollowsViewer      bool   `json:"follows_viewer"`
		FullName           string `json:"full_name"`
		HasChannel         bool   `json:"has_channel"`
		HasBlockedViewer   bool   `json:"has_blocked_viewer"`
		HighlightReelCount int    `json:"highlight_reel_count"`
		HasRequestedViewer bool   `json:"has_requested_viewer"`
		ID                 string `json:"id"`
		IsPrivate          bool   `json:"is_private"`
		IsVerified         bool   `json:"is_verified"`
		MutualFollowers    struct {
			AdditionalCount int      `json:"additional_count"`
			Usernames       []string `json:"usernames"`
		} `json:"mutual_followers"`
		ProfilePicURL                string      `json:"profile_pic_url"`
		ProfilePicURLHd              string      `json:"profile_pic_url_hd"`
		RequestedByViewer            bool        `json:"requested_by_viewer"`
		Username                     string      `json:"username"`
		ConnectedFbPage              interface{} `json:"connected_fb_page"`
		EdgeFelixCombinedPostUploads struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_felix_combined_post_uploads"`
		EdgeFelixCombinedDraftUploads struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_felix_combined_draft_uploads"`
		EdgeFelixVideoTimeline struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_felix_video_timeline"`
		EdgeFelixDrafts struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_felix_drafts"`
		EdgeFelixPendingPostUploads struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_felix_pending_post_uploads"`
		EdgeFelixPendingDraftUploads struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_felix_pending_draft_uploads"`
		EdgeOwnerToTimelineMedia struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool   `json:"has_next_page"`
				EndCursor   string `json:"end_cursor"`
			} `json:"page_info"`
			Edges []struct {
				Node struct {
					Typename           string `json:"__typename"`
					ID                 string `json:"id"`
					EdgeMediaToCaption struct {
						Edges []struct {
							Node struct {
								Text string `json:"text"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_media_to_caption"`
					Shortcode          string `json:"shortcode"`
					EdgeMediaToComment struct {
						Count int `json:"count"`
					} `json:"edge_media_to_comment"`
					CommentsDisabled bool `json:"comments_disabled"`
					TakenAtTimestamp int  `json:"taken_at_timestamp"`
					Dimensions       struct {
						Height int `json:"height"`
						Width  int `json:"width"`
					} `json:"dimensions"`
					DisplayURL  string `json:"display_url"`
					EdgeLikedBy struct {
						Count int `json:"count"`
					} `json:"edge_liked_by"`
					EdgeMediaPreviewLike struct {
						Count int `json:"count"`
					} `json:"edge_media_preview_like"`
					GatingInfo   interface{} `json:"gating_info"`
					MediaPreview string      `json:"media_preview"`
					Owner        struct {
						ID string `json:"id"`
					} `json:"owner"`
					ThumbnailSrc       string `json:"thumbnail_src"`
					ThumbnailResources []struct {
						Src          string `json:"src"`
						ConfigWidth  int    `json:"config_width"`
						ConfigHeight int    `json:"config_height"`
					} `json:"thumbnail_resources"`
					IsVideo        bool `json:"is_video"`
					VideoViewCount int  `json:"video_view_count"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"edge_owner_to_timeline_media"`
		EdgeSavedMedia struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_saved_media"`
		EdgeMediaCollections struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasNextPage bool        `json:"has_next_page"`
				EndCursor   interface{} `json:"end_cursor"`
			} `json:"page_info"`
			Edges []interface{} `json:"edges"`
		} `json:"edge_media_collections"`
	} `json:"user"`
}

type userData struct {
	Username        string
	ID              string
	Description     string
	Followers       int
	Following       int
	EngagementRatio float32
	AverageLikes    int
	AverageComments int
	Media           []mediaData
	AudienceTags    []Tag
	PublishedTags   []Tag
	FollowingTags   []Tag
	Email           string
	Private         bool
}

type Tag struct {
	Tag   string
	Ratio float64
	Count int
}

type LikerJson struct {
	Users []struct {
		Pk              int64  `json:"pk"`
		Username        string `json:"username"`
		FullName        string `json:"full_name"`
		IsPrivate       bool   `json:"is_private"`
		IsVerified      bool   `json:"is_verified"`
		ProfilePicURL   string `json:"profile_pic_url"`
		ProfilePicID    string `json:"profile_pic_id,omitempty"`
		LatestReelMedia int    `json:"latest_reel_media,omitempty"`
	} `json:"users"`
	UserCount int    `json:"user_count"`
	Status    string `json:"status"`
}

type mediaData struct {
	ID        string
	Code      string
	Likes     int
	Comments  int
	Caption   string
	Tags      []Tag
	LikerData likerData
}

type likerData struct {
	Username string
	fullName string
	picURL   string
}

type likerDatas []likerData

type endpoints struct {
	profileHttp string
	profilePage string
	likersPage  string
}

type Auth struct {
	Username       string
	Password       string
	gis            string
	client         *http.Client
	userAgent      string
	CsrfToken      string
	InstagramAjax  string
	InstagramGIS   string
	RequestedWith  string
	origin         string
	referer        string
	contentType    string
	sessionCookies string
	proxies        []string
}

const URL_Login = "https://instagram.com/accounts/login/ajax/"
const URL_Base = "https://instagram.com"
const userName = "user1"
const passWord = "password1"
const emailRegex = `\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3}`
const hashtagRegex = `#([a-zA-Z]+)`
const requestDelay = 50
const audienceLimit = 500

func Endpoints(username string, userID string) endpoints {
	endpoints := endpoints{}
	endpoints.likersPage = fmt.Sprintf("https://i.instagram.com/api/v1/media/%s/likers/", userID)
	endpoints.profilePage = fmt.Sprintf("https://www.instagram.com/%s/?__a=1", userName)
	fmt.Println(endpoints.profilePage)
	endpoints.profileHttp = fmt.Sprintf("https://www.instagram.com/%s", userName)
	return endpoints
}

func getExistingUsers() []string {
	var users []string
	//sQuery := elastic.NewMatchAllQuery()
	//tQuery := elastic.NewExistsQuery("Username")
	fmt.Println("Checking users in elastic index")
	svc := elasticClient.Scroll(elasticIndexName).Type(elasticTypeName).Size(10000)
	pages := 0
	docs := 0
	testIndexName := elasticIndexName

	for {
		res, err := svc.Do(context.TODO())
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
		}
		if res == nil {
			log.Print("expected results != nil; got nil")
		}
		if res.Hits == nil {
			log.Print("expected results.Hits != nil; got nil")
		}
		//if want, have := int64(3), res.Hits.TotalHits; want != have {
		//	log.Printf("expected results.Hits.TotalHits = %d; got %d", want, have)
		//}
		//if want, have := 1, len(res.Hits.Hits); want != have {
		//	log.Printf("expected len(results.Hits.Hits) = %d; got %d", want, have)
		//}

		pages++

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				log.Printf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := userData{}
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				log.Print(err)
			}
			//fmt.Print(item.Username)
			users = append(users, strings.TrimSpace(strings.ToLower(item.Username)))
			docs++
		}

		if len(res.ScrollId) == 0 {
			log.Printf("expected scrollId in results; got %q", res.ScrollId)
		}

		fmt.Print(len(users))
		pages++

	}
	return users
}

func checkSingleUserExists(username string) bool {
	check := false
	sQuery := elastic.NewTermQuery("Username", username)
	fmt.Println("Checking users in elastic index")
	res, err := elasticClient.Search(elasticIndexName).
		Index(elasticIndexName).
	//Type(elasticTypeName).
		Query(sQuery).
		//Size(10000).
	//Sort("Username", true).
		Do(context.Background())
	if err != nil {
		//panic(err)
		log.Print(err)
		return false
		// return err
	}
	var userDat userData
	for _, match := range res.Each(reflect.TypeOf(userDat)) {
		if t, ok := match.(userData); ok {
			if strings.TrimSpace(strings.ToLower(t.Username)) == username {
				check = true
				break
			}
		}
	}
	return check
}

func checkAlreadyScraped(input string, existing []string) bool {
	var check bool

	for _, s := range existing {
		if strings.ToLower(strings.TrimSpace(input)) == strings.ToLower(strings.TrimSpace(s)) {
			check = true
			//fmt.Printf("Already scraped %s", input)
			break
		}
		if input != s {
			check = false
		}
	}

	return check
}

func createScrapeQueue(input []string, existing []string) []string {
	var scrapeQueue []string
	exclude := map[string]bool{}
	for _, e := range existing {
		exclude[strings.ToLower(strings.TrimSpace(e))] = true
	}
	for _, s := range input {
		if !exclude[s] {
			scrapeQueue = append(scrapeQueue, s)
		}
	}
	return scrapeQueue
}

func cleanUserData(rawJson ProfileJson) userData {
	var tlikes int
	var tcomments int
	media := mediaData{}
	userData := userData{}
	userData.Username = rawJson.Graphql.User.Username
	userData.ID = rawJson.Graphql.User.ID
	userData.Description = rawJson.Graphql.User.Biography
	userData.Followers = rawJson.Graphql.User.EdgeFollowedBy.Count
	userData.Following = rawJson.Graphql.User.EdgeFollow.Count
	userData.AverageLikes = 0
	userData.AverageComments = 0
	userData.EngagementRatio = 0
	//fmt.Println(userData.Username)
	fmt.Println(len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges))
	for l := 0; l < len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges); l++ {
		ilikes := rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeLikedBy.Count
		tlikes += ilikes
		icomments := rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToComment.Count
		tcomments += icomments
		media.ID = rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.ID
		media.Likes = rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeLikedBy.Count
		if len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToCaption.Edges) > 0 {
			media.Caption = rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToCaption.Edges[0].Node.Text
		}
		userData.Media = append(userData.Media, media)
		//here
	}
	//fmt.Print(userData.Media)
	if len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges) > 0 {
		userData.AverageLikes = tlikes / len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges)
		userData.AverageComments = tcomments / len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges)
		userData.EngagementRatio = ((float32(userData.AverageLikes) + float32(userData.AverageComments)) / float32(userData.Followers)) * 100
	}
	//fmt.Println(userData.Media)
	return userData
}

func cleanUserDataLiker(rawJson ProfileJson) userData {
	var tlikes int
	var tcomments int
	media := mediaData{}
	userData := userData{}
	userData.Username = rawJson.Graphql.User.Username
	userData.ID = rawJson.Graphql.User.ID
	userData.Description = rawJson.Graphql.User.Biography
	userData.Followers = rawJson.Graphql.User.EdgeFollowedBy.Count
	userData.Following = rawJson.Graphql.User.EdgeFollow.Count
	userData.AverageLikes = 0
	userData.AverageComments = 0
	userData.EngagementRatio = 0
	//fmt.Println(userData.Username)
	fmt.Println(len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges))
	for l := 0; l < len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges); l++ {
		ilikes := rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeLikedBy.Count
		tlikes += ilikes
		icomments := rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToComment.Count
		tcomments += icomments
		media.ID = rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.ID
		media.Likes = rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeLikedBy.Count
		if len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToCaption.Edges) > 0 {
			media.Caption = rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToCaption.Edges[0].Node.Text
		}
		userData.Media = append(userData.Media, media)
	}
	if len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges) > 0 {
		userData.AverageLikes = tlikes / len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges)
		userData.AverageComments = tcomments / len(rawJson.Graphql.User.EdgeOwnerToTimelineMedia.Edges)
		userData.EngagementRatio = ((float32(userData.AverageLikes) + float32(userData.AverageComments)) / float32(userData.Followers)) * 100
	}
	//fmt.Println(userData.Media)
	return userData
}

func cleanUserDataShared(rawJson sharedData) userData {
	var tlikes int
	var tcomments int
	media := mediaData{}
	userData := userData{}
	userData.Username = rawJson.EntryData.ProfilePage[0].Graphql.User.Username
	userData.ID = rawJson.EntryData.ProfilePage[0].Graphql.User.ID
	userData.Description = rawJson.EntryData.ProfilePage[0].Graphql.User.Biography
	userData.Followers = rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeFollowedBy.Count
	userData.Following = rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeFollow.Count
	userData.AverageLikes = 0
	userData.AverageComments = 0
	userData.EngagementRatio = 0
	//fmt.Println(userData.Username)
	fmt.Println(len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges))
	if len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges) <= 0 {
		return userData
	}
	for l := 0; l < len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges); l++ {
		ilikes := rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeLikedBy.Count
		tlikes += ilikes
		icomments := rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToComment.Count
		tcomments += icomments
		media.ID = rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.ID
		media.Likes = rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeLikedBy.Count
		if len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToCaption.Edges) > 0 {
			media.Caption = rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges[l].Node.EdgeMediaToCaption.Edges[0].Node.Text
		}
		userData.Media = append(userData.Media, media)
	}
	if len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges) > 0 {
		userData.AverageLikes = tlikes / len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges)
		userData.AverageComments = tcomments / len(rawJson.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges)
		userData.EngagementRatio = ((float32(userData.AverageLikes) + float32(userData.AverageComments)) / float32(userData.Followers)) * 100
	}
	//fmt.Println(userData.Media)
	return userData
}

func cleanLikers(rawLikerJson LikerJson) likerDatas {
	//var likerData likerData
	var likerDats likerDatas
	var likerDat likerData
	//fmt.Println(rawLikerJson.Users[0].Username)
	for l := 0; l < len(rawLikerJson.Users); l++ {
		likerDat.Username = rawLikerJson.Users[l].Username
		likerDat.fullName = rawLikerJson.Users[l].FullName
		likerDat.picURL = rawLikerJson.Users[l].ProfilePicURL
		likerDats = append(likerDats, likerDat)
		//fmt.Println("Cleaned: ",likerDats[l].Username)
	}
	return likerDats
}

const (
	elasticIndexName = "instagram_profiles"
	elasticTypeName  = "instagram_profile"
	serverIP = "173.249.46.92"
)

var (
	elasticClient *elastic.Client
)

// var wg sync.WaitGroup

func InitElastic() {
	mapping := `{
		"mappings":{
			"` + elasticTypeName + `":{
				"properties": {
					"AverageComments": {
						"type": "long"
					},
					"AverageLikes": {
						"type": "long"
					},
					"Description": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"Email": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"EngagementRatio": {
						"type": "long"
					},
					"Followers": {
						"type": "long"
					},
					"Following": {
						"type": "long"
					},
					"ID": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"Private": {
						"type": "boolean"
					},
					"PublishedTags": {
						"type": "nested",
						"properties": {
							"Count": {
								"type": "long"
							},
							"Ratio": {
								"type": "float"
							},
							"Tag": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 256
									}
								}
							}
						}
					},
					"Username": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					}
				}
			}
		}
	}
}`

	ctx := context.Background()
	fmt.Println("Deleting Index..")
	deleteIndex, err := elasticClient.DeleteIndex(elasticIndexName).Do(ctx)
	if err != nil {
		// Handle error
		fmt.Printf("%v", err)
	}
	if !deleteIndex.Acknowledged {
		fmt.Printf("%s", "Index doesn't exist")
		// Not acknowledged
	}

	fmt.Println("Creating Index..")
	_, err = elasticClient.CreateIndex(elasticIndexName).BodyString(mapping).Do(ctx)
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func main() {
	var err error
	for {
		elasticClient, err = elastic.NewClient(
			//elastic.SetURL(fmt.Sprintf("http://%s:9200"), serverIP),
			elastic.SetURL("http://173.249.46.92:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	var auth Auth
	auth = loadProxies("proxies.txt", auth)
	randP := rand.Intn(len(auth.proxies) - 1)
	auth = setProxy(auth.proxies[randP], auth)

	cookies, err := ioutil.ReadFile("cookie.txt")
	if err != nil {
		log.Println(err)
		fmt.Println("no cookies found, logging in..")
		auth = login(userName, passWord, auth)
	}
	if string(cookies) != "" {
		auth.sessionCookies = string(cookies)
		fmt.Println("session cookies found")
	}

	users := readUsernameList("users.txt")
	existing := getExistingUsers()

	scrapeQueue := createScrapeQueue(users, existing)

	//log.Print(len(users))
	//log.Print(len(existing))
	log.Print(len(scrapeQueue))

	for _, username := range scrapeQueue {
		scrapedUser := scrapeSingleUserShared(username, auth)
		if scrapedUser.Private == true || scrapedUser.ID == "" || scrapedUser.Username == "" || len(scrapedUser.Media) <= 0 {
			fmt.Print("error:",scrapedUser.Username)
			continue
		}
		fmt.Print("scraping")
		scrapedUser.AudienceTags = scrapeAudienceTags(scrapedUser, auth)
		scrapedUser.Media = []mediaData{}
		saveElasticUser(scrapedUser)
	}
}

func checkSingleUserExistsMultithread(user <-chan string, check chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	checkExist := ""
	username := <-user
	username = strings.TrimSpace(strings.ToLower(username))
	sQuery := elastic.NewTermQuery("Username",username)
	//fmt.Println("Checking users in elastic index")
	res, err := elasticClient.Search(elasticIndexName).
		Index(elasticIndexName).
		Query(sQuery).
		Do(context.Background())
	if err != nil {
		check <- checkExist
		return
	}
	var userDat userData
	//fmt.Print("Checking user exists:", username)
	for _, match := range res.Each(reflect.TypeOf(userDat)) {
		if t, ok := match.(userData); ok {
			if strings.TrimSpace(strings.ToLower(t.Username)) == strings.TrimSpace(strings.ToLower(username)) {
				checkExist = username
				break
			}
		}
	}
	check <- checkExist
	return
}

func checkWorker(done chan bool) {
	end := <-done
	for {
		if end == true {
			return
		}
	}
}

func main221() {
	var err error

	elasticClient, err = elastic.NewClient(
		elastic.SetURL("http://173.249.46.92:9200"),
		elastic.SetSniff(false),
		elastic.SetHttpClient(&http.Client{Transport: &http.Transport{
			//MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: true,
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
		}}),

	)
	if err != nil {
		log.Print(err)
		time.Sleep(3 * time.Second)
	}

	//InitElastic()

	var auth Auth

	auth = loadProxies("proxies.txt", auth)
	log.Print(auth.proxies[0])
	randP := rand.Intn(len(auth.proxies) - 1)
	auth = setProxy(auth.proxies[randP], auth)

	cookies, err := ioutil.ReadFile("cookie.txt")
	if err != nil {
		log.Println(err)
		fmt.Println("no cookies found, logging in..")
		auth = login(userName, passWord, auth)
	}
	if string(cookies) != "" {
		auth.sessionCookies = string(cookies)
		fmt.Println("session cookies found")
	}

	var users []string
	//var existing []string
	var scrapedUser userData
	//var alreadyScraped bool

	users = readUsernameList("users.txt")
	wg := sync.WaitGroup{}
	userC := make(chan string)
	checkC := make(chan string)
	counter := 0
	workers := 100
	done := make(chan bool)
	var scrapeQueue []string
	//existing = getExistingUsers()
	for range users {
		wg.Add(1)
		go checkSingleUserExistsMultithread(userC, checkC, &wg)
	}
	for _, username := range users {
		log.Print(username)
		username = strings.TrimSpace(strings.ToLower(username))
	}
	for range users {
		checkUser := <- checkC
		if checkUser == "" {
			continue
		}
		scrapeQueue = append(scrapeQueue, checkUser)
	}

	for range users {
		if counter == len(users) {
			for u:=0; u<workers; u++ {
				done <- true
				break
			}
		}
	}

	wg.Wait()
	for _, username := range scrapeQueue {
		scrapedUser = scrapeSingleUserShared(username, auth)
		if scrapedUser.Private == true || scrapedUser.ID == "" || scrapedUser.Username == "" || len(scrapedUser.Media) <= 0 {
			fmt.Print("error")
			continue
		}
		fmt.Print("scraping")
		scrapedUser.AudienceTags = scrapeAudienceTags(scrapedUser, auth)
		scrapedUser.Media = []mediaData{}
		saveElasticUser(scrapedUser)
	}
}

func getSharedData(auth Auth) *sharedData {
	var shared sharedData
	//strings.index

	return &shared
}
func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func generateGis(gisCode string, params string) string {
	var gisHash string

	gis := fmt.Sprintf("%s:%s", gisCode, params)
	fmt.Println(gis)
	h := md5.New()
	h.Write([]byte(gis))
	gisHash = fmt.Sprintf("%x", h.Sum(nil))
	return gisHash
}

func login(username string, password string, auth Auth) Auth {
	var mid string
	var csrftoken string

	println("Initiate...")
	req, _ := http.NewRequest("GET", URL_Base, nil)
	client := &http.Client{}
	req.Header.Set("cookie", "ig_cb=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "csrftoken" {
			csrftoken = cookies[i].Value
		}
		if cookies[i].Name == "mid" {
			mid = cookies[i].Value
		}
	}

	fmt.Printf("CSRF Token : %s\n", csrftoken)
	fmt.Printf("mid : %s\n", mid)

	println("Logging in...")
	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)
	req, _ = http.NewRequest("POST", URL_Login, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
	req.Header.Set("x-csrftoken", csrftoken)
	req.Header.Set("cookie", fmt.Sprintf("csrftoken=%s; mid=%s;", csrftoken, mid))
	req.Header.Set("referer", URL_Base)

	resp, err = client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", string(body))

	cookies = resp.Cookies()

	var sessionCsrf string
	var sessionCookie string

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "csrftoken" && cookies[i].Value != "" {
			sessionCsrf = cookies[i].Value
		}

		sessionCookie += fmt.Sprintf("%s=%s; ", cookies[i].Name, cookies[i].Value)
	}

	println("Get Account Information...")
	req, _ = http.NewRequest("GET", "https://www.instagram.com/accounts/edit/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	req.Header.Set("cookie", sessionCookie)
	req.Header.Set("x-csrftoken", sessionCsrf)
	req.Header.Set("referer", URL_Base+"/")

	resp, err = client.Do(req)
	body, _ = ioutil.ReadAll(resp.Body)

	var shared sharedData
	window := between(string(body), `window._sharedData = `, `;</script>`)
	_ = json.Unmarshal([]byte(window), &shared)

	fmt.Printf("Shared Data : %+v", shared)
	fmt.Println(sessionCookie)
	auth.sessionCookies = sessionCookie
	auth.Username = username
	auth.Password = password
	auth.CsrfToken = sessionCsrf
	auth.userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36"
	ioutil.WriteFile("cookie.txt", []byte(sessionCookie), 0777)
	return auth
}

func scrapeUsername(uName string, auth Auth) userData { //, useProxy bool) userData {
	var userScraped ProfileJson
	var userCleaned userData
	var tags []Tag
	var tagsTrimmed []Tag
	var email string

	randP := rand.Intn(len(auth.proxies) - 1)
	auth = setProxy(auth.proxies[randP], auth)
	urL, err := url.Parse(fmt.Sprintf("https://instagram.com/%s/?__a=1", uName))
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	req, err := http.NewRequest("GET", urL.String(), nil)
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
	req.Header.Set("x-csrftoken", auth.CsrfToken)
	req.Header.Set("referer", URL_Base)
	req.Header.Set("cookie", auth.sessionCookies)
	//resp, err := http.DefaultClient.Do(req)
	resp, err := auth.client.Do(req)
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	//	}
	defer resp.Body.Close()
	req.Close = true
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Close = true
	fmt.Println("Scraping:" + uName)
	tags = scrapeTags(string(bytes))
	fmt.Println(uName, "tags:", tags)
	uErr := json.Unmarshal(bytes, &userScraped)
	if userScraped.Graphql.User.IsPrivate == true || userScraped.Graphql.User.ID == "" || uErr != nil {
		//skip user
		log.Print(uErr)
		userCleaned.Private = true
		//fmt.Println("User is private:", userScraped.Graphql.User.Username)
		return userCleaned
	}
	fmt.Println("Cleaning user:", userScraped.Graphql.User.Username)
	userCleaned = cleanUserData(userScraped)
	email = scrapeEmail(userCleaned.Description)
	ByField(tags, "Ratio")
	tags = reverse(tags)
	tagsMax := 10
	if len(tags) < tagsMax {
		tagsMax = len(tags)
	}
	for ta := 0; ta < tagsMax; ta++ {
		tagsTrimmed = append(tagsTrimmed, tags[ta])
	}
	userCleaned.PublishedTags = tagsTrimmed
	userCleaned.Email = email

	return userCleaned
}

func scrapeSingleUserShared(uName string, auth Auth) userData { //, useProxy bool) userData {
	var userScraped sharedData
	var userCleaned userData
	var tags []Tag
	var email string
	randP := rand.Intn(len(auth.proxies) - 1)
	auth = setProxy(auth.proxies[randP], auth)
	urL, err := url.Parse(fmt.Sprintf("https://instagram.com/%s/", uName))
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	req, err := http.NewRequest("GET", urL.String(), nil)
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	resp, err := auth.client.Do(req)
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return userCleaned
	}
	resp.Close = true
	req.Close = true
	window := between(string(bytes), `window._sharedData = `, `;</script>`)
	uErr := json.Unmarshal([]byte(window), &userScraped)
	//log.Print("Profile page:", len(userScraped.EntryData.ProfilePage))
	if len(userScraped.EntryData.ProfilePage) <= 0 {
		//log.Print("User has no media:", uName)
		return userCleaned
	}
	userCleaned.Private = userScraped.EntryData.ProfilePage[0].Graphql.User.IsPrivate
	userCleaned.ID = userScraped.EntryData.ProfilePage[0].Graphql.User.ID
	userCleaned.Username = userScraped.EntryData.ProfilePage[0].Graphql.User.Username
	if userCleaned.Private == true || userCleaned.ID == "" || uErr != nil {
		return userCleaned
	}
	fmt.Println("Scraping:" + uName)
	tags = scrapeTags(window)
	fmt.Println(uName, "Published tags:", tags)
	fmt.Println("Cleaning user:", userCleaned.Username)
	userCleaned = cleanUserDataShared(userScraped)
	email = scrapeEmail(userCleaned.Description)
	ByField(tags, "Ratio")
	tags = reverse(tags)
	userCleaned.PublishedTags = tags
	userCleaned.Email = email
	var tagsTrimmed []Tag
	tagsMax := 10
	if len(tags) < tagsMax {
		tagsMax = len(tags)
	}
	for ta := 0; ta < tagsMax; ta++ {
		tagsTrimmed = append(tagsTrimmed, tags[ta])
	}
	userCleaned.PublishedTags = tagsTrimmed
	return userCleaned
}

func scrapeUsernameLiker(uChan <-chan string, auth Auth, liker chan<- userData, wg *sync.WaitGroup) { //, useProxy bool) userData {
		defer wg.Done()
		var userScraped ProfileJson
		var userCleaned userData
		var tags []Tag
		var email string
		var uName string
		uName = <-uChan
		//time.Sleep(time.Millisecond * 50)
		randP := rand.Intn(len(auth.proxies) - 1)
		auth = setProxy(auth.proxies[randP], auth)
		urL, err := url.Parse(fmt.Sprintf("https://instagram.com/%s/?__a=1", uName))
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		req, err := http.NewRequest("GET", urL.String(), nil)
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
		req.Header.Set("x-csrftoken", auth.CsrfToken)
		req.Header.Set("referer", URL_Base)
		req.Header.Set("cookie", auth.sessionCookies)
		resp, err := auth.client.Do(req)
		//resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		//	}
		defer resp.Body.Close()
		req.Close = true
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		resp.Close = true
		fmt.Println("Scraping:" + uName)
		tags = scrapeTags(string(bytes))
		fmt.Println(uName, "tags:", tags)
		uErr := json.Unmarshal(bytes, &userScraped)
		if userScraped.Graphql.User.IsPrivate == true || userScraped.Graphql.User.ID == "" || uErr != nil {
			//skip user
			userCleaned.Private = true
			liker <- userCleaned
			//fmt.Println("User is private:", userScraped.Graphql.User.Username)
			return
		}
		fmt.Println("Cleaning user:", userScraped.Graphql.User.Username)
		userCleaned = cleanUserDataLiker(userScraped)
		email = scrapeEmail(userCleaned.Description)
		ByField(tags, "Ratio")
		tags = reverse(tags)
		userCleaned.PublishedTags = tags
		userCleaned.Email = email
		liker <- userCleaned

		var tagsTrimmed []Tag
		if email != "" {
			tagsMax := 10
			if len(tags) < tagsMax {
				tagsMax = len(tags)
			}
			for ta := 0; ta < tagsMax; ta++ {
				tagsTrimmed = append(tagsTrimmed, tags[ta])
			}
			userCleaned.Media = []mediaData{}
			userCleaned.PublishedTags = tagsTrimmed
			saveElasticUser(userCleaned)
		}
		return
}

func scrapeUsernameShared(uChan <-chan string, auth Auth, liker chan<- userData, wg *sync.WaitGroup) { //, useProxy bool) userData {
		defer wg.Done()
		var userScraped sharedData
		var userCleaned userData
		var tags []Tag
		var email string
		uName := <-uChan
		randP := rand.Intn(len(auth.proxies) - 1)
		auth = setProxy(auth.proxies[randP], auth)
		urL, err := url.Parse(fmt.Sprintf("https://instagram.com/%s/", uName))
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		req, err := http.NewRequest("GET", urL.String(), nil)
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		resp, err := auth.client.Do(req)
		if err != nil {
			log.Print(err)
			liker <- userCleaned
			return
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			liker <- userCleaned
			return
		}
		resp.Close = true
		req.Close = true
		window := between(string(bytes), `window._sharedData = `, `;</script>`)
		uErr := json.Unmarshal([]byte(window), &userScraped)
		if len(userScraped.EntryData.ProfilePage) <= 0 {
			liker <- userCleaned
			return
		}
		userCleaned.Private = userScraped.EntryData.ProfilePage[0].Graphql.User.IsPrivate
		userCleaned.ID = userScraped.EntryData.ProfilePage[0].Graphql.User.ID
		userCleaned.Username = userScraped.EntryData.ProfilePage[0].Graphql.User.Username
		if userCleaned.Private == true || userCleaned.ID == "" || uErr != nil {
			liker <- userCleaned
			return
		}
		fmt.Println("Scraping:" + uName)
		tags = scrapeTags(window)
		fmt.Println(uName, "Published tags:", tags)
		fmt.Println("Cleaning user:", userCleaned.Username)
		userCleaned = cleanUserDataShared(userScraped)
		email = scrapeEmail(userCleaned.Description)
		ByField(tags, "Ratio")
		tags = reverse(tags)
		userCleaned.PublishedTags = tags
		userCleaned.Email = email
		liker <- userCleaned
	var tagsTrimmed []Tag
	if email != "" {
		tagsMax := 10
		if len(tags) < tagsMax {
			tagsMax = len(tags)
		}
		for ta := 0; ta < tagsMax; ta++ {
			tagsTrimmed = append(tagsTrimmed, tags[ta])
		}
		userCleaned.Media = []mediaData{}
		userCleaned.PublishedTags = tagsTrimmed
		saveElasticUser(userCleaned)
	}
	return
}

func scrapeEmail(text string) string {
	var email string
	r := regexp.MustCompile(emailRegex)
	email = r.FindString(text)
	fmt.Println(email)
	return email
}

func readUsernameList(filepath string) []string {
	fileContent, _ := ioutil.ReadFile(filepath)
	lines := strings.Split(string(fileContent), "\n")
	return lines
}

func scrapeUsernameList(filePath string, auth Auth) { //useProxy bool) []userData {
	// var user []userData
	//maybe scrape input slice of users instead of filepath directly
	fileContent, _ := ioutil.ReadFile(filePath)
	fmt.Println(string(fileContent))
	lines := strings.Split(string(fileContent), "\n")

	for i := 0; i < len(lines); i++ { //for each line in file, scrape use
		go scrapeUsername(lines[i], auth) //scrape user from file at index i
	}

	// return user
}

func scrapeAudienceTags(uData userData, auth Auth) []Tag {
	var likesCounter []int
	var audienceTags []Tag
	var audienceTagsTrimmed []Tag
	tagsCount := 0
	ag := make(map[string]Tag)
	fmt.Print("Scraping audience tags:", uData.ID)
	if len(uData.Media) <= 0 {
		return nil
	}
	for _, media := range uData.Media {
		likesCounter = append(likesCounter, media.Likes)
	}
	sort.Ints(likesCounter)
	mostLikedIndex := len(uData.Media) - 1
	likers := scrapeLikers(uData.Media[mostLikedIndex].ID, auth)
	if likers == nil || len(likers) <= 0 {
		return nil
	}
	maxThreads := len(likers)
	scrapedChan := make(chan userData, maxThreads)
	likerChan := make(chan string, maxThreads)
	wg := sync.WaitGroup{}
	for i := 0; i < len(likers); i++ {
		wg.Add(1)
		go scrapeUsernameShared(likerChan, auth, scrapedChan, &wg)
	}
	for _, liker := range likers {
		//time.Sleep(time.Millisecond* 50)
		fmt.Println(uData.ID,"Scraping tags from:", liker.Username)
		likerChan <- liker.Username
	}
	for a := 0; a < len(likers); a++ {
		scrapedUser := <- scrapedChan
		if scrapedUser.Username == "" || scrapedUser.Private==true || scrapedUser.ID == "" {
			continue
		}
		publishedTags := scrapedUser.PublishedTags
		tagsCount += len(publishedTags)
		for _, v := range publishedTags {
			c, ok := ag[v.Tag]
			if !ok {
				c.Tag = v.Tag
				c.Count = 0
			}
			c.Count++
			ag[v.Tag] = c
		}
	}
	wg.Wait()
	close(likerChan)
	close(scrapedChan)
	for _, k := range ag {
		k.Ratio = float64(k.Count) / float64(tagsCount)
		audienceTags = append(audienceTags, k)
	}
	ByField(audienceTags, "Ratio")
	audienceTags = reverse(audienceTags)
	tagsMax := 15
	if len(audienceTags) < tagsMax {
		tagsMax = len(audienceTags)
	}
	for ta := 0; ta < tagsMax; ta++ {
		audienceTagsTrimmed = append(audienceTagsTrimmed, audienceTags[ta])
	}
	return audienceTagsTrimmed
}

func scrapeFollowingTags(uData userData, auth Auth) []Tag {
	var tags []Tag
	//var followingJson FollowingJson

	getFollowingUser(uData, auth)

	//tags := scrapeTags
	//go to following JSON and scrape all users for tags

	return tags
}

func getFollowingUser(uData userData, auth Auth) { //followers struct
	loop := true
	afterQuery := ""
	after := ""
	page := 1
	totalfollowing := 0
	for loop == true {
		fmt.Println("following page:", page)
		if after != "" {
			afterQuery = fmt.Sprintf(", \"after\": \"%s\"", after)
		}
		nextpage := fmt.Sprintf("{\"id\":\"%s\",\"first\":\"%d\"%s}", uData.ID, 50, afterQuery)

		nextpage = strings.Replace(nextpage, " ", "", -1)
		qs := fmt.Sprintf("query_hash=9335e35a1b280f082a47b98c5aa10fa4&variables=%s", nextpage)

		url := fmt.Sprintf("https://www.instagram.com/graphql/query/?%s", qs)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		req.Header.Set("cookie", auth.sessionCookies)
		req.Header.Set("x-csrftoken", auth.CsrfToken)
		req.Header.Set("referer", URL_Base)
		// req.Header.Set("x-instagram-gis", gisHash)

		resp, err := auth.client.Do(req)
		if err != nil {
			fmt.Printf("length: %v", err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var folData FollowingJson
		json.Unmarshal(body, &folData)
		if len(folData.Data.User.EdgeFollow.Edges) <= 0 {
			fmt.Println(string(body))
			break
		}

		for _, edges := range folData.Data.User.EdgeFollow.Edges {
			fmt.Println(uData.Username, "is following:", edges.Node.Username)
		}
		fmt.Println(len(folData.Data.User.EdgeFollow.Edges))
		totalfollowing += len(folData.Data.User.EdgeFollow.Edges)
		if folData.Data.User.EdgeFollow.PageInfo.HasNextPage == true {
			//proxy := randomizeProxy()
			//auth = setProxy(proxy.address, auth)
			randP := rand.Intn(len(auth.proxies) - 1)
			auth = setProxy(auth.proxies[randP], auth)
			after = folData.Data.User.EdgeFollow.PageInfo.EndCursor
			page++
			time.Sleep(time.Second * requestDelay)
		}
		if folData.Data.User.EdgeFollow.PageInfo.HasNextPage == false {
			loop = false
		}
	}
	fmt.Println("Total following:", totalfollowing)
	return

}

func scrapeLikers(mID string, auth Auth) likerDatas {
	var likersComplete LikerJson
	randP := rand.Intn(len(auth.proxies) - 1)
	auth = setProxy(auth.proxies[randP], auth)
	//parseLikersUrl := "https://i.instagram.com/api/v1/media/"+mID+"/likers/"
	parseLikersUrl := fmt.Sprintf("https://i.instagram.com/api/v1/media/%s/likers/", mID)
	req, err := http.NewRequest("GET", parseLikersUrl, nil)
	if err != nil {
		log.Print(err)
		 return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
	req.Header.Set("x-csrftoken", auth.CsrfToken)
	req.Header.Set("referer", URL_Base)
	req.Header.Set("cookie", auth.sessionCookies)

	response, err := auth.client.Do(req)
	//response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(bytes, &likersComplete)
	req.Close = true
	response.Close = true
	likersCleaned := cleanLikers(likersComplete)
	fmt.Println(len(likersCleaned))
	return likersCleaned
}

func scrapeTags(text string) []Tag {
	var tags []Tag
	var tagi Tag
	var lowertag string

	m := make(map[string]Tag)

	var pattern = regexp.MustCompile(hashtagRegex)
	matches := pattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		//fmt.Println("Found hashtag:", match[1])
		//if strings.Contains(match[1], `\`) == true {
		//	break
		//}
		lowertag = strings.ToLower(match[1])
		c, ok := m[lowertag]
		if !ok {
			c.Tag = lowertag
			c.Count = 0
		}
		c.Count++
		m[lowertag] = c
	}

	for _, v := range m {
		tagi.Tag = v.Tag
		tagi.Count = v.Count
		tagi.Ratio = float64(v.Count) / float64(len(matches))
		tags = append(tags, tagi)
	}

	return tags
}

func loadProxies(filePath string, auth Auth) Auth {
	proxyList, _ := ioutil.ReadFile(filePath)
	lines := strings.Split(string(proxyList), "\n")
	fmt.Print(len(lines))
	auth.proxies = lines
	return auth
}

func setProxy(proxy string, auth Auth) Auth {
	proxy = strings.TrimSpace(proxy)
	proxy = fmt.Sprintf("http://%s",proxy)
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			//MaxIdleConns:       10,
			IdleConnTimeout:    12 * time.Second,
			DisableCompression: true,
			Dial: (&net.Dialer{
				Timeout: 12 * time.Second,
				KeepAlive: 12 * time.Second,
			}).Dial,
		},
	}
	auth.client = client
	return auth
}

func saveElasticUser(user userData) error {
	fmt.Println("Saving to elastic")
	id := shortid.MustGenerate()
	_, err := elasticClient.Index().
		Index(elasticIndexName).
		Type(elasticTypeName).
		Id(id).
		BodyJson(user).
		Do(context.Background())
	if err != nil {
		//panic(err)
		return nil
	}

	return nil

	//	name := user.Name
	//	bio := user.Description
}

func ByField(data interface{}, field string) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("ByField takes a slice as data")
	}
	if v.Len() == 0 {
		return
	}
	t := v.Type().Elem()
	if t.Kind() != reflect.Struct {
		panic("ByField takes a slice of structs as data")
	}
	fieldStruct, ok := t.FieldByName(field)
	if !ok {
		panic("ByField cannot find field " + field)
	}
	v1 := v.Index(0)
	f1 := v1.FieldByIndex(fieldStruct.Index)
	var helper sort.Interface
	switch f1.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		helper = ByFieldHelperInt{
			v:         v,
			size:      v.Len(),
			tmp:       reflect.New(v.Type().Elem()),
			keyValues: makeIntKeys(v, fieldStruct.Index, v.Len()),
		}
	case reflect.Float32, reflect.Float64:
		helper = ByFieldHelperFloat{
			v:         v,
			size:      v.Len(),
			tmp:       reflect.New(v.Type().Elem()),
			keyValues: makeFloatKeys(v, fieldStruct.Index, v.Len()),
		}
	case reflect.String:
		helper = ByFieldHelperString{
			v:         v,
			size:      v.Len(),
			tmp:       reflect.New(v.Type().Elem()),
			keyValues: makeStringKeys(v, fieldStruct.Index, v.Len()),
		}
	default:
		panic("Cannot compare " + f1.Type().String() + " values")
	}
	sort.Sort(helper)
}

func swapValues(slice reflect.Value, tmp reflect.Value, i, j int) {
	vi := slice.Index(i)
	vj := slice.Index(j)
	tmp.Elem().Set(vi)
	vi.Set(vj)
	vj.Set(tmp.Elem())
}

type ByFieldHelperInt struct {
	v         reflect.Value
	size      int
	tmp       reflect.Value
	keyValues []int64
}

func (t ByFieldHelperInt) Len() int { return t.size }

func (t ByFieldHelperInt) Swap(i, j int) {
	swapValues(t.v, t.tmp, i, j)
	t.keyValues[i], t.keyValues[j] = t.keyValues[j], t.keyValues[i]
}

func (t ByFieldHelperInt) Less(i, j int) bool {
	return t.keyValues[i] < t.keyValues[j]
}

func makeIntKeys(v reflect.Value, fieldId []int, len int) []int64 {
	keys := make([]int64, len)
	for i := 0; i < len; i++ {
		vi := v.Index(i)
		fi := vi.FieldByIndex(fieldId)
		keys[i] = fi.Int()
	}
	return keys
}

type ByFieldHelperString struct {
	v         reflect.Value
	size      int
	tmp       reflect.Value
	keyValues []string
}

func (t ByFieldHelperString) Len() int { return t.size }

func (t ByFieldHelperString) Swap(i, j int) {
	swapValues(t.v, t.tmp, i, j)
	t.keyValues[i], t.keyValues[j] = t.keyValues[j], t.keyValues[i]
}

func (t ByFieldHelperString) Less(i, j int) bool {
	return t.keyValues[i] < t.keyValues[j]
}

func makeStringKeys(v reflect.Value, fieldId []int, len int) []string {
	keys := make([]string, len)
	for i := 0; i < len; i++ {
		vi := v.Index(i)
		fi := vi.FieldByIndex(fieldId)
		keys[i] = fi.String()
	}
	return keys
}

type ByFieldHelperFloat struct {
	v         reflect.Value
	size      int
	tmp       reflect.Value
	keyValues []float64
}

func (t ByFieldHelperFloat) Len() int { return t.size }

func (t ByFieldHelperFloat) Swap(i, j int) {
	swapValues(t.v, t.tmp, i, j)
	t.keyValues[i], t.keyValues[j] = t.keyValues[j], t.keyValues[i]
}

func (t ByFieldHelperFloat) Less(i, j int) bool {
	return t.keyValues[i] < t.keyValues[j]
}

func makeFloatKeys(v reflect.Value, fieldId []int, len int) []float64 {
	keys := make([]float64, len)
	for i := 0; i < len; i++ {
		vi := v.Index(i)
		fi := vi.FieldByIndex(fieldId)
		keys[i] = fi.Float()
	}
	return keys
}

func reverse(s []Tag) []Tag {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
