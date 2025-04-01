import { LogOut, Trophy, Share2 } from "lucide-react";
import { supabase } from "../lib/supabase";
import { Player } from "../types";
import toast from "react-hot-toast";

interface NavbarProps {
  player: Player | null;
  onShowLeaderboard: () => void;
  onLogout: () => void;
  isGuest?: boolean;
  guestStats?: { correct: number; total: number };
  onShowAuthForm?: () => void;
}

export function Navbar({
  player,
  onShowLeaderboard,
  onLogout,
  isGuest = false,
  guestStats,
  onShowAuthForm,
}: NavbarProps) {
  if (isGuest) {
    // Single row layout for guest mode
    return (
      <nav className="bg-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between py-4">
            <div className="flex items-center">
              <img
                src="https://api.dicebear.com/7.x/pixel-art/svg?seed=%22guest%22"
                alt="Guest"
                className="w-8 h-8 rounded-full"
              />
              <span className="ml-3 font-medium text-gray-900">Guest</span>
              <div className="ml-6 flex items-center text-sm text-gray-600">
                <span>
                  Correct: {guestStats ? guestStats.correct : 0}/
                  {guestStats ? guestStats.total : 0}
                </span>
              </div>
            </div>
            <div className="flex items-center">
              <button
                onClick={onShowAuthForm}
                className="px-3 py-2 rounded-md text-sm font-medium text-blue-500 hover:bg-gray-100"
              >
                Login / Register
              </button>
            </div>
          </div>
        </div>
      </nav>
    );
  }

  if (!player) {
    return (
      <nav className="bg-white shadow-md p-4">
        <div className="text-center text-gray-700">Loading profile...</div>
      </nav>
    );
  }

  const handleShare = async () => {
    const {
      data: { session },
    } = await supabase.auth.getSession();
    const inviteLink = `${window.location.origin}/?invite-code=${session?.user.id}`;

    if (navigator.share) {
      try {
        await navigator.share({
          title: "Join me on Globetrotter!",
          text: "Can you beat my score in this amazing travel guessing game?",
          url: inviteLink,
        });
      } catch (error: any) {
        if (error.name !== "AbortError") {
          console.error("Error sharing:", error);
        }
      }
    } else {
      await navigator.clipboard.writeText(inviteLink);
      toast.success("Invite link copied to clipboard!");
    }
  };

  // Logged-in mode: Original two-layer layout
  return (
    <nav className="bg-white shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row md:justify-between py-4 md:h-16 md:py-0">
          <div className="flex items-center justify-center md:justify-start">
            <img
              src={player.avatar}
              alt={player.name}
              className="w-8 h-8 rounded-full"
            />
            <span className="ml-3 font-medium text-gray-900">
              {player.name}
            </span>
            <div className="ml-6 flex items-center space-x-2 text-sm text-gray-600">
              <span>Score: {Math.round(player.score * 100)}</span>
              <span>|</span>
              <span>
                Correct: {player.correct_answers}/{player.total_attempts}
              </span>
            </div>
          </div>

          <div className="flex items-center justify-center md:justify-end space-x-4 mt-4 md:mt-0">
            <button
              onClick={handleShare}
              className="flex items-center px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100"
            >
              <Share2 className="w-5 h-5 mr-2" />
              <span className="whitespace-nowrap">Invite Friends</span>
            </button>
            <button
              onClick={onShowLeaderboard}
              className="flex items-center px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100"
            >
              <Trophy className="w-5 h-5 mr-2" />
              <span className="whitespace-nowrap">Leaderboard</span>
            </button>
            <button
              onClick={onLogout}
              className="flex items-center px-3 py-2 rounded-md text-sm font-medium text-red-600 hover:bg-red-50"
            >
              <LogOut className="w-5 h-5 mr-2" />
              <span className="whitespace-nowrap">Logout</span>
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}
