import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { useEffect, useState } from "react";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const useImageSlider = (interval = 5000) => {
  const images: string[] = Array.from(
    { length: 7 },
    (_, i) => `scenes/${i + 1}.png`
  );

  useEffect(() => {
    images.forEach((src) => {
      const img = new Image();
      img.src = src;
    });
  }, [images]);

  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    const sliderInterval = setInterval(() => {
      setTimeout(() => {
        setCurrentIndex((prevIndex) => (prevIndex + 1) % images.length);
      }, 500);
    }, interval);

    return () => clearInterval(sliderInterval);
  }, [images, interval]);

  return images[currentIndex];
};
