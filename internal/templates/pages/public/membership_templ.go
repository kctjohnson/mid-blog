// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.598
package public

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/kctjohnson/mid-blog/internal/templates/components"
import "github.com/kctjohnson/mid-blog/internal/db/models"

func Membership(user *models.User) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			templ_7745c5c3_Err = components.Navbar(user).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"w-full flex flex-col items-center py-8 px-4\"><div class=\"prose\"><h1>Membership</h1><h3>For Anonymous Users:</h3><p>As explorers on the surface of our satirical seas, you get to dip your toes in with a sneak peek at the opening paragraph of any article. It's just enough to tickle your funny bone and spark your interest. While enjoying this glimpse of what our AI has on offer, please note that interacting through comments, as well as the ability to rate posts, are special features reserved for our registered members. To fully immerse in the discourse and hilarity, step out of the shadows of anonymity and sign up!</p><h3>For Registered Users (Free Members):</h3><p>We extend a hearty welcome to the inner circle, where your status as a registered user grants you unfettered access to the complete array of our free articles. Not only can you devour every word of these pieces, but you are also empowered to rate and comment on them, adding your voice to the symphony of our community's conversations and critiques. Engage, deliberate, and be an active participant in the ever-growing narrative of our platform.</p><h3>For Premium Users:</h3><p>For the aficionados who yearn for the richest flavors of our satirical feast, premium access awaits. As a subscriber, not only do you enjoy all the privileges of free membership, but you are also privy to the full breadth of our premium content — no more teasers, just complete, unadulterated AI-crafted articles that promise to both challenge and entertain. Engage in lively debate, offer your ratings, leave thought-provoking comments, and help curate the experience for all by sharing your insights on these exclusive pieces.</p><p>Your participation as a paid subscriber is not only about enjoying the depths of satire but also about sustaining the creativity and continuity of <strong>Midblog</strong>. Your support allows this unique space to thrive, ever-expanding the boundaries of humor with each AI-generated word.</p><h3>Step into the Limelight of Laughter and Commentary</h3><p>Don't just watch from the wings; join us center stage by signing up or subscribing to become a pivotal part of our satirical endeavor. Ready for incisive wit and in-depth discussion? Your membership grants you access to the full suite of features our blog platform offers. Register for free or become a premium subscriber today, and never miss a beat (or a joke) with <strong>Midblog</strong>.</p></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = components.Layout("Membership").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}