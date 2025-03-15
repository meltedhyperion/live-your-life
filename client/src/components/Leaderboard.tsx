import React from 'react';
import { X } from 'lucide-react';
import { LeaderboardEntry } from '../types';

interface LeaderboardProps {
    leaderboard: LeaderboardEntry[];
    onClose: () => void;
}

export function Leaderboard({ leaderboard, onClose }: LeaderboardProps) {
    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl">
                <div className="flex items-center justify-between p-6 border-b">
                    <h2 className="text-2xl font-bold">Leaderboard</h2>
                    <button
                        onClick={onClose}
                        className="text-gray-500 hover:text-gray-700"
                    >
                        <X className="w-6 h-6" />
                    </button>
                </div>

                <div className="p-6">
                    <div className="space-y-4">
                        {leaderboard.map((entry, index) => (
                            <div
                                key={index}
                                className="flex items-center p-4 bg-gray-50 rounded-lg"
                            >
                                <div className="flex items-center flex-1">
                                    <span className="font-bold text-lg w-8">{index + 1}</span>
                                    <img
                                        src={entry.avatar}
                                        alt={entry.name}
                                        className="w-10 h-10 rounded-full"
                                    />
                                    <div className="ml-4">
                                        <div className="font-medium">{entry.name}</div>
                                        <div className="text-sm text-gray-500">
                                            {entry.correct_answers}/{entry.total_attempts} correct
                                        </div>
                                    </div>
                                </div>
                                <div className="text-2xl font-bold text-blue-600">
                                    {Math.round(entry.score * 100)}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}