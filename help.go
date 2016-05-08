package main


var helpString = `<b>Ingress Glyph Bot</b>

This bot makes stickers of Ingress glyphs and glyph sequences.

Open a conversation with the bot, and try sending it messages like this:

► One glyph: <code>Portal</code>
► Multiple glyphs: <code>Clear All, Open All, Discover, Truth</code>
► Un-used glyphs: <code>N'Zeer, Victory</code>

Or you can add the bot as a member of your chat, and use it like this:

► <code>/glyphs Breathe, Inside, XM, Lose, Self</code>

<b>Troubleshooting</b>

This is <b>beta software</b>. If the bot isn't responding, it's probably crashed. Message @liaungyip and complain!

► The glyph names must be spelled exactly like in the scanner / Glyphtionary.
    ✓ <code>Portal, Human</code>
    ✗ <code>Portals, Humans</code> (Extra <code>s</code>)

► The glyph names must be separated using commas. If you forget the commas, it won't work.
    ✓ <code>Help, Enlightened, Capture, All, Portal</code>
    ✗ <code>Help Enlightened Capture All Portal</code>

Bugs? Improvements? <a href="http://github.com/LiaungYip/tgGlyphBot">Github: LiaungYip/tgGlyphBot</a>`