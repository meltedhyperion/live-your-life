interface HeadoutButtonProps {
  city: string
}

export default function HeadoutButton({ city }: HeadoutButtonProps) {
  const cityName = city.split(',')[0];
  const formattedCity = cityName.toLowerCase();
  const searchUrl = `https://www.headout.com/search/?q=${formattedCity}&c=`

  return (
    <a
      href={searchUrl}
      target="_blank"
      rel="noopener noreferrer"
      className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-blue-600 to-purple-600 text-white font-medium hover:from-blue-700 hover:to-purple-700 transition-all shadow-md hover:shadow-lg"
    >
      <svg
        width="24"
        height="24"
        viewBox="0 0 223 120"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="w-10 h-10"
      >
        <path
          d="M213.991 49.0276C198.086 78.8969 158.712 93.7881 118.122 93.7881C85.8755 93.7881 54.1507 89.521 26.5111 77.2424C92.5681 84.2961 156.539 81.5094 213.991 49.0276ZM213.991 36.5747C194.696 30.8273 177.312 29.1727 159.321 29.3469C113.689 29.6952 67.536 43.4543 26.7718 61.4804C59.1919 32.656 107.431 11.6691 155.757 11.6691C178.79 11.6691 205.3 21.5094 213.991 36.5747ZM222.857 43.1059C222.857 12.4528 176.704 0 150.368 0C109.082 0 74.4024 15.0653 44.7637 31.611C42.9384 28.3019 40.3309 27.6923 36.2458 27.1698C30.1616 26.6473 26.4241 26.3861 21.8175 26.3861C16.7763 26.3861 9.21454 26.7344 3.30417 27.6923C0.435898 28.2148 -0.520191 30.479 0.262064 33.9623C2.86958 43.4543 6.78085 51.9884 11.909 60C9.64912 67.8374 10.4314 75.2395 16.9502 81.3353C14.6903 90.479 13.995 100.493 14.3426 109.55C14.8642 113.295 16.7763 114.34 19.9053 113.817C31.5523 111.988 42.8515 107.896 52.6731 102.671C55.1068 101.364 56.2367 99.8839 56.9321 98.0551L85.7017 102.496C90.1345 107.286 94.5672 111.379 99.9561 115.907C103.52 118.52 107.344 120 112.038 120C132.985 119.826 150.629 116.691 165.231 111.118C170.967 109.289 175.226 101.626 181.571 90.9144C195.652 83.7736 222.857 70.3628 222.857 43.1059Z"
          fill="white"
        />
      </svg>
      <span className="leading-none">to {cityName}</span>
    </a>
  )
}
