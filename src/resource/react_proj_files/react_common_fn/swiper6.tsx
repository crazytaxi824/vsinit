// this is swiper6,
// swiper7: https://swiperjs.com/migration-guide

import React from 'react';
import SwiperCore, { Autoplay, Pagination, Navigation } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/react';

// 引入 css 文件
import 'swiper/swiper.min.css';
import 'swiper/components/pagination/pagination.min.css';
import 'swiper/components/navigation/navigation.min.css';

// 使用 pagination 和 自动播放
SwiperCore.use([Pagination, Autoplay, Navigation]);

export function SwiperImg(): JSX.Element {
  return (
    <Swiper
      style={{
        width: '98%',
        height: '300px',
        borderRadius: '0.4rem',
        borderStyle: 'dashed',
        borderWidth: '1px',
      }}
      navigation
      slidesPerView={1} // 每页显示1个图片
      autoplay={{
        delay: 2500, // 自动播放 2.5s
        disableOnInteraction: false, // 可以手动滑动
      }}
      loop // 可以循环无限滚动, 可以左右滚动
      pagination={{
        clickable: true, // slide 下面的小点点可以点击翻页
      }}
    >
      <SwiperSlide>
        <img
          alt="1.png"
          src="http://localhost:18080/img/1.png"
          style={{
            width: '98%',
            height: '300px',
            objectFit: 'cover',
            objectPosition: '0 0',
          }}
        />
      </SwiperSlide>
      <SwiperSlide>
        <img
          alt="1.png"
          src="http://localhost:18080/img/1.png"
          style={{
            width: '98%',
            height: '300px',
            objectFit: 'cover',
            objectPosition: '0 0',
          }}
        />
      </SwiperSlide>
      <SwiperSlide>
        <img
          alt="1.png"
          src="http://localhost:18080/img/1.png"
          style={{
            width: '98%',
            height: '300px',
            objectFit: 'cover',
            objectPosition: '0 0',
          }}
        />
      </SwiperSlide>
      {/* <SwiperSlide>Slide 4</SwiperSlide>
      <SwiperSlide>Slide 5</SwiperSlide> */}
    </Swiper>
  );
}
