package cleaner

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Deware-PK/goniverse-cli/internal/config"
	"golang.org/x/net/html"
)

type HTMLCleaner struct {
	Cfg *config.Config
}

func NewHTMLCleaner(cfg *config.Config) *HTMLCleaner {
	return &HTMLCleaner{Cfg: cfg}
}

func (c *HTMLCleaner) Clean(input string) string {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return input
	}

	var buf bytes.Buffer
	c.traverse(doc, &buf)

	return buf.String()
}

func (c *HTMLCleaner) traverse(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.ElementNode && c.Cfg.AllowedTags[n.Data] {
		buf.WriteString("<" + n.Data)

		for _, attr := range n.Attr {
			key := attr.Key
			val := attr.Val

			if key == "src" || key == "href" || key == "alt" || key == "title" {
				if n.Data == "img" && key == "src" {
					_, filename := filepath.Split(val)
					val = "images/" + filename
				}

				buf.WriteString(fmt.Sprintf(` %s="%s"`, key, val))
			}
		}

		buf.WriteString(">")
	} else if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.traverse(child, buf)
	}

	if n.Type == html.ElementNode && c.Cfg.AllowedTags[n.Data] {
		if n.Data != "img" && n.Data != "br" && n.Data != "hr" {
			buf.WriteString("</" + n.Data + ">")
		}
	}
}

func (c *HTMLCleaner) WrapHTML(title string, content string, currentChapter int, totalChapters int) string {

	prevLink := `<span class="nav-btn disabled" style="opacity:0.5; cursor:default;">&larr; Previous</span>`
	if currentChapter > 1 {
		prevLink = fmt.Sprintf(`<a href='chapter_%03d.html' class="nav-btn">&larr; Previous</a>`, currentChapter-1)
	}

	nextLink := `<span class="nav-btn disabled" style="opacity:0.5; cursor:default;">Next &rarr;</span>`
	if currentChapter < totalChapters {
		nextLink = fmt.Sprintf(`<a href='chapter_%03d.html' class="nav-btn">Next &rarr;</a>`, currentChapter+1)
	}

	headAndScript := `
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Crimson+Pro:ital,wght@0,400;0,600;0,700;1,400;1,600&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">

    <style>
        /* =========================================
           1. CSS Variables & Themes
           ========================================= */
        :root {
            /* Base Dimensions */
            --spacing-unit: 1rem;
            --container-width: 760px;
            --header-height: 60px;
            
            /* Typography Scale */
            --font-base-size: 18px;
            --font-heading: 'Crimson Pro', serif;
            --font-body: 'Crimson Pro', serif;
            --font-ui: 'Inter', sans-serif;
            --line-height-body: 1.75;
            --letter-spacing-body: 0.01em;

            /* Transition */
            --transition-speed: 0.3s;
        }

        /* Theme: Light */
        [data-theme="light"] {
            --bg-body: #F9F9F7;
            --text-main: #2C2C2C;
            --text-secondary: #595959;
            --accent-color: #3B82F6;
            --quote-bg: #F0F0EE;
            --quote-border: #D1D5DB;
            --ui-bg: rgba(255, 255, 255, 0.95);
            --ui-border: #E5E7EB;
            --ui-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            --img-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
        }

        /* Theme: Sepia */
        [data-theme="sepia"] {
            --bg-body: #F4ECD8;
            --text-main: #433422;
            --text-secondary: #746656;
            --accent-color: #A05A2C;
            --quote-bg: #EAE0CC;
            --quote-border: #C8BCA8;
            --ui-bg: rgba(244, 236, 216, 0.95);
            --ui-border: #DBCDB4;
            --ui-shadow: 0 4px 20px rgba(67, 52, 34, 0.1);
            --img-shadow: 0 10px 15px -3px rgba(67, 52, 34, 0.15);
        }

        /* Theme: Dark */
        [data-theme="dark"] {
            --bg-body: #1A1A1A;
            --text-main: #E5E5E5;
            --text-secondary: #A3A3A3;
            --accent-color: #60A5FA;
            --quote-bg: #262626;
            --quote-border: #404040;
            --ui-bg: rgba(30, 30, 30, 0.95);
            --ui-border: #333333;
            --ui-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
            --img-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
        }

        /* =========================================
           2. Global Styles
           ========================================= */
        * { box-sizing: border-box; margin: 0; padding: 0; }

        body {
            background-color: var(--bg-body);
            color: var(--text-main);
            font-family: var(--font-body);
            font-size: var(--font-base-size);
            line-height: var(--line-height-body);
            letter-spacing: var(--letter-spacing-body);
            transition: background-color var(--transition-speed), color var(--transition-speed);
            -webkit-font-smoothing: antialiased;
            padding-bottom: 120px;
            cursor: pointer; /* Hint for toggle UI */
        }

        /* Top Nav */
        .top-nav {
            position: fixed;
            top: 0; left: 0; right: 0;
            height: var(--header-height);
            display: flex;
            align-items: center;
            justify-content: center;
            background: var(--bg-body);
            z-index: 100;
            transition: transform var(--transition-speed), opacity var(--transition-speed), background-color var(--transition-speed);
            border-bottom: 1px solid var(--ui-border);
        }
        
        .top-nav.ui-hidden { transform: translateY(-100%); opacity: 0; pointer-events: none; }
        
        .top-nav span {
            font-family: var(--font-ui);
            font-size: 0.9rem;
            color: var(--text-secondary);
            font-weight: 600;
            letter-spacing: 0.05em;
            text-transform: uppercase;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 90%;
        }

        /* Container */
        .container {
            max-width: var(--container-width);
            margin: calc(var(--header-height) + 2rem) auto 0;
            padding: 0 24px;
            cursor: auto;
        }

        /* Content Styling */
        h1, h2, h3, h4, h5, h6 {
            font-family: var(--font-heading);
            color: var(--text-main);
            margin-top: 2em; margin-bottom: 1em;
            font-weight: 700; line-height: 1.3;
        }
        h1 { font-size: 2.2em; text-align: center; margin-top: 1em; }
        p { margin-bottom: 1.5em; text-align: justify; text-justify: inter-word; }
        b, strong { font-weight: 700; color: var(--text-main); }
        i, em { font-style: italic; }
        
        blockquote {
            margin: 2em 0;
            padding: 1.5rem 2rem;
            background-color: var(--quote-bg);
            border-left: 4px solid var(--accent-color);
            border-radius: 0 12px 12px 0;
            color: var(--text-secondary);
            font-style: italic;
        }

        img {
            display: block; max-width: 100%; height: auto;
            margin: 2.5em auto; border-radius: 12px;
            box-shadow: var(--img-shadow);
        }

        /* Navigation Links */
        .nav-links {
            display: flex; justify-content: space-between;
            margin-top: 4rem; padding-top: 2rem;
            border-top: 1px solid var(--quote-border);
            font-family: var(--font-ui);
            gap: 10px;
        }

        .nav-btn {
            display: inline-flex; align-items: center;
            padding: 12px 24px; border-radius: 50px;
            background-color: transparent;
            border: 1px solid var(--quote-border);
            color: var(--text-main);
            text-decoration: none; font-weight: 500;
            transition: all 0.2s ease; cursor: pointer;
            font-size: 0.9rem;
        }
        .nav-btn:hover:not(.disabled) {
            background-color: var(--quote-bg);
            border-color: var(--text-secondary);
            transform: translateY(-2px);
        }

        /* Floating Toolbar */
        .toolbar {
            position: fixed; bottom: 30px; left: 50%;
            transform: translateX(-50%);
            background-color: var(--ui-bg);
            backdrop-filter: blur(12px);
            padding: 12px 24px; border-radius: 100px;
            display: flex; gap: 24px; align-items: center;
            box-shadow: var(--ui-shadow);
            border: 1px solid var(--ui-border);
            z-index: 1000; font-family: var(--font-ui);
            transition: transform var(--transition-speed), opacity var(--transition-speed);
        }

        .toolbar.ui-hidden { transform: translate(-50%, 150%); opacity: 0; pointer-events: none; }

        .control-group { display: flex; align-items: center; gap: 12px; }
        .divider { width: 1px; height: 24px; background-color: var(--ui-border); }

        .theme-selector {
            width: 24px; height: 24px; border-radius: 50%;
            border: 2px solid var(--ui-border); cursor: pointer;
            transition: transform 0.2s;
        }
        .theme-selector:hover { transform: scale(1.2); }
        .theme-selector.active { border-color: var(--accent-color); transform: scale(1.1); }

        .tool-btn {
            background: none; border: none; cursor: pointer;
            width: 36px; height: 36px; border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            transition: all 0.2s; color: var(--text-secondary);
            font-weight: 600; font-family: serif;
        }
        .tool-btn:hover { background-color: var(--quote-bg); color: var(--text-main); }
        .font-small { font-size: 14px; }
        .font-large { font-size: 22px; }

        @media (max-width: 600px) {
            .toolbar { width: 90%; justify-content: space-between; padding: 12px 16px; bottom: 20px; }
            :root { --font-base-size: 16px; }
            h1 { font-size: 1.8em; }
        }
    </style>

    <script>
        const htmlEl = document.documentElement;
        let currentFontSize = 18;

        function setTheme(themeName) {
            htmlEl.setAttribute('data-theme', themeName);
            localStorage.setItem("goniverse_theme", themeName);
            
            document.querySelectorAll('.theme-selector').forEach(el => {
                el.classList.remove('active');
                if(el.getAttribute('onclick').includes(themeName)) {
                    el.classList.add('active');
                }
            });
        }

        function adjustFont(change) {
            currentFontSize += change;
            if (currentFontSize < 14) currentFontSize = 14;
            if (currentFontSize > 26) currentFontSize = 26;
            htmlEl.style.setProperty('--font-base-size', currentFontSize + 'px');
            localStorage.setItem("goniverse_font", currentFontSize);
        }

        // Toggle UI when clicking background
        document.addEventListener('click', (e) => {
            const selection = window.getSelection();
            if (selection.toString().length > 0) return;
            if (e.target.closest('.toolbar') || e.target.closest('.nav-btn')) return;

            const toolbar = document.querySelector('.toolbar');
            const topNav = document.querySelector('.top-nav');
            
            toolbar.classList.toggle('ui-hidden');
            topNav.classList.toggle('ui-hidden');
        });

        // Init
        (function() {
            const savedTheme = localStorage.getItem("goniverse_theme") || 'light';
            setTheme(savedTheme);
            
            const savedFont = parseInt(localStorage.getItem("goniverse_font"));
            if (savedFont) {
                currentFontSize = savedFont;
                htmlEl.style.setProperty('--font-base-size', currentFontSize + 'px');
            }
        })();
    </script>`

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="th" data-theme="light">
<head>
    <meta charset="utf-8">
    <title>%s</title>
    %s
</head>
<body>

    <div class="top-nav">
        <span>%s</span>
    </div>

    <div class="container" id="book-content">
        <h1>%s</h1>
        %s
    </div>

    <div class="container">
        <div class="nav-links">
            %s
            <a href='index.html' class="nav-btn">List of contents</a>
            %s
        </div>
    </div>

    <div class="toolbar">
        <div class="control-group">
            <div class="theme-selector" style="background: #F9F9F7;" onclick="setTheme('light')" title="Light"></div>
            <div class="theme-selector" style="background: #F4ECD8;" onclick="setTheme('sepia')" title="Sepia"></div>
            <div class="theme-selector" style="background: #333333;" onclick="setTheme('dark')" title="Dark"></div>
        </div>
        <div class="divider"></div>
        <div class="control-group">
            <button class="tool-btn font-small" onclick="adjustFont(-2)">A</button>
            <button class="tool-btn font-large" onclick="adjustFont(2)">A</button>
        </div>
    </div>

</body>
</html>`, title, headAndScript, title, title, content, prevLink, nextLink)
}
