import { useEffect, useState } from "react";
import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  Users,
  ListMusic,
  Search,
  Volume2,
  Settings,
} from "lucide-react";
import "./App.css";

type Track = {
  id: string;
  title: string;
  artist: string;
  duration: number;
};

const demoQueue: Track[] = [
  {
    id: "1",
    title: "Auryvex Demo",
    artist: "Auryvex",
    duration: 214,
  },
];

function formatTime(seconds: number) {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);

  return `${mins}:${secs.toString().padStart(2, "0")}`;
}

function App() {
  const [playing, setPlaying] = useState(false);
  const [position, setPosition] = useState(0);
  const [volume, setVolume] = useState(80);
  const [showQueue, setShowQueue] = useState(false);
  const [search, setSearch] = useState("");

  const currentTrack = demoQueue[0];

  useEffect(() => {
    if (!playing) return;

    const timer = window.setInterval(() => {
      setPosition((value) => {
        if (value >= currentTrack.duration) {
          setPlaying(false);
          return 0;
        }

        return value + 1;
      });
    }, 1000);

    return () => window.clearInterval(timer);
  }, [playing, currentTrack.duration]);

  return (
    <main className="app">
      <header className="topbar">
        <div className="brand">
          <div className="logo">A</div>
          <div>
            <h1>Auryvex</h1>
            <span>Shared Room</span>
          </div>
        </div>

        <div className="actions">
          <button title="Members">
            <Users size={20} />
          </button>

          <button
            title="Queue"
            onClick={() => setShowQueue((value) => !value)}
          >
            <ListMusic size={20} />
          </button>

          <button title="Settings">
            <Settings size={20} />
          </button>
        </div>
      </header>

      <section className="room">
        <div className="cover">
          <div className="cover-letter">A</div>
        </div>

        <div className="track-info">
          <span className="label">NOW PLAYING</span>
          <h2>{currentTrack.title}</h2>
          <p>{currentTrack.artist}</p>
        </div>

        <div className="progress-area">
          <input
            type="range"
            min="0"
            max={currentTrack.duration}
            value={position}
            onChange={(event) =>
              setPosition(Number(event.target.value))
            }
          />

          <div className="times">
            <span>{formatTime(position)}</span>
            <span>{formatTime(currentTrack.duration)}</span>
          </div>
        </div>

        <div className="controls">
          <button>
            <SkipBack size={22} />
          </button>

          <button
            className="play"
            onClick={() => setPlaying((value) => !value)}
          >
            {playing ? <Pause size={26} /> : <Play size={26} />}
          </button>

          <button>
            <SkipForward size={22} />
          </button>
        </div>

        <div className="volume">
          <Volume2 size={18} />

          <input
            type="range"
            min="0"
            max="100"
            value={volume}
            onChange={(event) => setVolume(Number(event.target.value))}
          />
        </div>
      </section>

      <section className="search">
        <Search size={19} />

        <input
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Search music or paste a YouTube link..."
        />
      </section>

      {showQueue && (
        <section className="queue">
          <div className="queue-header">
            <h3>Queue</h3>
            <span>{demoQueue.length} track</span>
          </div>

          {demoQueue.map((track) => (
            <div className="queue-item" key={track.id}>
              <div className="queue-cover">A</div>

              <div>
                <strong>{track.title}</strong>
                <span>{track.artist}</span>
              </div>

              <time>{formatTime(track.duration)}</time>
            </div>
          ))}
        </section>
      )}
    </main>
  );
}

export default App;
