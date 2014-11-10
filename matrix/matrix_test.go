package matrix

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/tokens"
)

func TestTree(t *testing.T) {
	docs := []string{
		`Billabong, Castelli, ***First time on The Clymb: Richer Poorer*** EARTH Watches, Kelty, urban footwear, mountain ukuleles, raft the Zambezi River, and more!`,
		`Zombie Alert: Gear up to Survive + Big Shot Commuter Bikes are Coming Up, and Much More. Take a Look!`,
		`Reach the rooftop of Africa, watch the sunrise over Machu Picchu, earn Level 1 Avalanche certification, and trek through Torres del Paine!`,
		`We're Totally Ruggin Out Over These New Prints!`,
		`Don’t Miss Your Chance: Register for the Twilio MMS Deep Dive`,
		`The Twilio Newsletter: Toll-Free SMS, Twilio CX for Chromebooks, and More`,
		`Andrew, save up to 40% on top travel deals.`,
		`Welcome Rewards® Account Summary`,
		`Free Tour to Rome for One Amazing Teenage African-American Female.`,
		`BIG UPDATE! Italy Food and Wine Tour`,
		`Today's Headlines: Lobbyists, Bearing Gifts, Pursue Attorneys General`,
		`Costco Services Update: College ID Theft`,
		`Win A Trip to The Lovies in London`,
		`The ONLY 5 dresses you need in your closet + Fashion week lessons that work IRL`,
		`Take notes...`,
		`Final hours: a secret so good, it's impossible not to share it`,
		`AWS: Recent News`,
		`Let’s hear it for the red, white & blue jeans (+ up to 50% off)`,
		`For the Dad you look up to…`,
		`One for me, one for you…`,
		`This is effortless sexy.`,
		`Win a Free Gift Card: Fill Out a Quick Survey to Enter`,
		`Want to Earn a $10 REI E-Gift Card?`,
		`In the Kitchen with Ina: The Barefoot Contessa Shares the Secrets to Her Success`,
		`LAST CHANCE: Get These Fall Pieces Before They're Gone For Good`,
		`Last Chance! Save 50% on ★★★★★ Ebooks and Training Videos`,
		`Save Up To 60% On Our AUTUMN SALE`,
		`Costco Services Health Insurance Marketplace`,
		`Our Stylist's Picks: October Essentials`,
		`Mid Month Drop: Wool-Blend Toggle Coat`,
		`Netflix optimizes constantly. You can too.`,
		`Take one step to improve your business`,
		`How Chrome Industries optimized retail for the web`,
		`72 people you never thought would A/B test`,
		`On Sale Today - 30% off QuickBooks! Ends soon!`,
		`Kyle, your improved rewards site is here!`,
		`Alert - 2 Days Only! Save 40% on QuickBooks`,
		`XFINITY® Deals Alert`,
		`5 Ways to Ignite Your Content Marketing With Presentations`,
		`TurboTax Notice: Your Privacy Statement`,
		`Are you taking advantage of all your benefits, Kyle?`,
		`Earn extra points when you shop online`,
		`Start Learning Tableau Today`,
		`EMC Finds Insight in Big Data, Fast`,
		`Tableau News: The Art of Analytics, Viz Roundup, and Much More`,
		`Create on Codecademy - Introducing Codebits`,
		`BOO! from Sprinkles`,
		`we're crushing on cobalt`,
		`last day to get your first outfit for $25!`,
		`Breakfast is free, Friday or Saturday only`,
		`See what's hot this week on Roku channels`,
		`Apple Pay available now`,
		`Dash your way to 20 Bonus Stars`,
		`Women Only: Your Outdoor Style Guide.`,
		`Enjoy one espresso, and get another to share`,
		`Afternoon tea and a bonus`,
		`Let's Autumn two-gether: Buy one, share one`,
		`Daily Deals`,
		`Imagine the travel: Win 1 million Expedia+ points`,
		`50% off hotels, just for you`,
		`Aveda Institute Portland, Invites You to Visit Us Soon`,
		`It's Almost Over - Your Free Minutes Expire Friday`,
		`NOVEMBER NEW PRODUCTS! Discover Your New Favorite!`,
		`Sweet Supplies + Halloween Savings!`,
		`BOO-k a flight right now!`,
		`Halfway through fall CSA today`,
		`Unmanned Antares Rocket Explodes Shortly After Takeoff`,
	}

	classes := [][]string{
		[]string{`Daily.Deal`},
		[]string{`Preview.Teaser`},
		[]string{`Discount`, `Money.Off`},
		[]string{`Product.Feature`},
		[]string{`Education`, `Rich.Content`, `Webinar`, `Event`},
		[]string{`Rich.Content`, `Newsletter`},
		[]string{`Rewards`, `Product.Usage`},
		[]string{`Rewards`, `Product.Usage`},
		[]string{`Giveaway.Contest`},
		[]string{`Education`, `Rich.Content`, `Video`},
		[]string{`Rich.Content`, `Newsletter`},
		[]string{`Product.Feature`, `Rich.Content`, `Newsletter`},
		[]string{`Giveaway.Contest`, `Event`},
		[]string{`Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Product.Feature`, `Sale`, `Discount`, `Percentage.Off`},
		[]string{`Sale`, `Discount`, `ActNow`},
		[]string{`Rich.Content`, `Blog`, `Newsletter`, `Video`, `Event`},
		[]string{`Discount`, `Percentage.Off`, `ActNow`},
		[]string{`Sale`, `Discount`, `Percentage.Off`, `Holiday.Seasonal`, `ActNow`},
		[]string{`Discount`, `Percentage.Off`, `Buy.One..Get.One`, `ActNow`},
		[]string{`Product.Feature`, `Discount`, `Percentage.Off`, `Free.Gift`, `Holiday.Seasonal`, `ActNow`},
		[]string{`Free.Gift`, `Survey`},
		[]string{`Discount`, `Money.Off`, `Referral`},
		[]string{`Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Product.Feature`, `ActNow`},
		[]string{`Discount`, `Percentage.Off`, `ActNow`},
		[]string{`Sale`, `Discount`, `Percentage.Off`, `Holiday.Seasonal`},
		[]string{`Product.Feature`, `Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Product.Feature`},
		[]string{`Product.Feature`},
		[]string{`Education`, `Rich.Content`, `Blog`},
		[]string{``},
		[]string{`Rich.Content`, `Blog`, `Video`},
		[]string{`Rich.Content`, `Blog`},
		[]string{`Discount`, `Percentage.Off`},
		[]string{`Product.Feature`},
		[]string{`Discount`, `Percentage.Off`, `ActNow`},
		[]string{`Discount`, `Education`},
		[]string{`Education`, `Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Notification.Alert`},
		[]string{`Product.Feature`, `Rich.Content`, `Newsletter`},
		[]string{`Product.Feature`, `Rewards`},
		[]string{`Product.Feature`, `Rich.Content`, `Blog`, `Newsletter`, `Video`},
		[]string{`Education`, `Rich.Content`, `Video`},
		[]string{`Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Product.Feature`, `Education`, `Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Holiday.Seasonal`},
		[]string{`Product.Feature`, `Discount`, `Money.Off`, `Giveaway.Contest`, `Thank.You`},
		[]string{`Discount`},
		[]string{`Free.Gift`, `ActNow`},
		[]string{`Product.Feature`, `Personal.Recommendations`},
		[]string{`Product.Feature`},
		[]string{`Rewards`, `Holiday.Seasonal`},
		[]string{`Discount`, `Percentage.Off`},
		[]string{`Discount`, `Buy.One..Get.One`},
		[]string{`Discount`, `Percentage.Off`, `Rewards`, `ActNow`},
		[]string{`Discount`, `Buy.One..Get.One`, `Holiday.Seasonal`},
		[]string{`Daily.Deal`},
		[]string{`Giveaway.Contest`, `Rich.Content`, `Newsletter`},
		[]string{`Discount`, `Percentage.Off`},
		[]string{`Reactivation`},
		[]string{`Free.Gift`, `ActNow`},
		[]string{`Product.Feature`, `Rich.Content`, `Blog`, `Newsletter`},
		[]string{`Sale`, `Discount`, `Percentage.Off`, `Holiday.Seasonal`},
		[]string{`Sale`, `Discount`, `Percentage.Off`, `Holiday.Seasonal`, `ActNow`},
		[]string{`Event`},
		[]string{`Notification.Alert`},
	}

	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	tree := NewTree()

	for i, doc := range docs {
		ngrams := tokenizer.Parse(doc)
		tree.Learn(ngrams, classes[i])
	}

	fmt.Printf("\n\n-------------------------------\n\n")
	for _, doc := range docs {
		fmt.Printf("\n\n\nDoc: %s\n", doc)
		predicted := tree.Predict(tokenizer.Parse(doc))

		for key, prob := range predicted {
			if prob > 0.2 {
				fmt.Printf("\n\t%s: %f", key, prob)
			}
		}
	}
}
