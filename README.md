## Info
Experimental TUI toolkit based on [tview](https://github.com/rivo/tview). Everything is a list! (buttons, inputs areas, text areas as options in lists, pages as lists) This is a failed experiment rather than something useful, TUI slop, if you will. On top of this sand castle there is a chat application, for real time interactions between users in-memory sharded map is used [xsync](https://github.com/puzpuzpuz/xsync), for event rate limiting used valkey and for database used postgres with default configuration. Backend behind nginx which is used as a rate limiting proxy and all of this is glued together with docker compose. Of course this is a toy example and there is even no proper authentication handling, but at least it's not broken and end-to-end test are fully passing, deterministic, and reproducible.


<p align="center" width="100%">
<video src="https://github.com/user-attachments/assets/eb0b21a5-89b7-4961-99c6-3b60f38030db" width="80%" controls></video>
</p>

### Run:
REDIS_DEBUG="true" docker compose up --build
