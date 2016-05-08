package main


var helpString = `<b>Ingress Glyph Bot</b>

This bot makes stickers of Ingress glyphs and glyph sequences.

Try sending it messages like this:

► One glyph: <code>Portal</code>
► Multiple glyphs: <code>Clear All, Open All, Discover, Truth</code>
► Un-used glyphs: <code>N'Zeer, Victory</code>

Or you can add it to your chat, and use it like this:

► <code>/glyphs Breathe, Inside, XM, Lose, Self</code>

<b>Troubleshooting</b>

► The glyph names must be spelled exactly like in the scanner / Glyphtionary.
    ✓ <code>Portal</code>
    ✗ <code>Portals</code> (Extra <code>s</code>)


► The glyph names must be separated using commas - not spaces, semicolons, dashes, or anything else.
    ✓ <code>Help, Enlightened, Capture, All, Portal</code>
    ✗ <code>Help Enlightened Capture All Portal</code>

Bugs? Improvements? <a href="http://github.com/LiaungYip/tgGlyphBot">Github: LiaungYip/tgGlyphBot</a>`