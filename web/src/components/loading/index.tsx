import { FC } from "react";
import { css, keyframes } from "@emotion/react";
import styled from "@emotion/styled";
import { useAntdTheme } from "@lepton-dashboard/hooks/use-antd-theme";

const Container = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const hideShow = keyframes`
  0% { opacity: 0; }
  50% { opacity: 1; }
  100% { opacity: 0; }
`;
export const Loading: FC = () => {
  const theme = useAntdTheme();
  return (
    <Container
      css={css`
        background: ${theme.colorBgContainer};
      `}
    >
      <svg width="100px" height="100px" x="0px" y="0px" viewBox="0 0 85 85">
        <path
          css={css`
            fill-rule: evenodd;
            clip-rule: evenodd;
            fill: #2d9cdb;
            animation: ${hideShow} 2s ease infinite;
          `}
          d="M75.9,48.1V36.9c0-2,0-3.1-0.1-3.9c0-0.4-0.1-0.6-0.1-0.7c-0.1-0.3-0.2-0.5-0.4-0.7c-0.1,0-0.2-0.2-0.6-0.4
	c-0.7-0.5-1.6-1-3.3-2l-9.7-5.6c-1.7-1-2.7-1.5-3.4-1.9c-0.4-0.2-0.6-0.3-0.6-0.3c-0.3-0.1-0.6-0.1-0.9,0c-0.1,0-0.3,0.1-0.6,0.3
	c-0.7,0.4-1.7,0.9-3.4,1.9l-9.7,5.6c-1.7,1-2.7,1.5-3.3,2c-0.3,0.2-0.5,0.4-0.6,0.4c-0.2,0.2-0.3,0.5-0.4,0.7c0,0.1,0,0.3-0.1,0.7
	c0,0.8-0.1,1.9-0.1,3.9v11.2c0,2,0,3.1,0.1,3.9c0,0.4,0.1,0.6,0.1,0.7c0.1,0.3,0.2,0.5,0.4,0.7c0.1,0,0.2,0.2,0.6,0.4
	c0.7,0.5,1.6,1,3.3,2l9.7,5.6c1.7,1,2.7,1.5,3.4,1.9c0.4,0.2,0.6,0.3,0.6,0.3c0.3,0.1,0.6,0.1,0.9,0c0.1,0,0.3-0.1,0.6-0.3
	c0.7-0.4,1.7-0.9,3.4-1.9l9.7-5.6c1.7-1,2.7-1.5,3.3-2c0.3-0.2,0.5-0.4,0.6-0.4c0.2-0.2,0.3-0.5,0.4-0.7c0-0.1,0-0.3,0.1-0.7
	C75.9,51.2,75.9,50.1,75.9,48.1z M75.7,52.7C75.7,52.7,75.7,52.7,75.7,52.7C75.7,52.7,75.7,52.7,75.7,52.7z M75.3,53.4
	C75.3,53.4,75.3,53.4,75.3,53.4C75.3,53.4,75.3,53.4,75.3,53.4z M57.7,63.7C57.7,63.7,57.7,63.7,57.7,63.7
	C57.7,63.7,57.7,63.7,57.7,63.7z M56.9,63.7C56.9,63.7,56.9,63.7,56.9,63.7C56.9,63.7,56.9,63.7,56.9,63.7z M39.3,53.4
	C39.3,53.4,39.3,53.4,39.3,53.4C39.3,53.4,39.3,53.4,39.3,53.4z M38.9,52.7C38.9,52.7,38.9,52.7,38.9,52.7
	C38.9,52.7,38.9,52.7,38.9,52.7z M38.9,32.3C38.9,32.3,38.9,32.3,38.9,32.3C38.9,32.3,38.9,32.3,38.9,32.3z M39.3,31.6
	C39.3,31.6,39.3,31.6,39.3,31.6C39.3,31.6,39.3,31.6,39.3,31.6z M56.9,21.4C56.9,21.4,56.9,21.4,56.9,21.4
	C56.9,21.4,56.9,21.4,56.9,21.4z M57.7,21.4C57.7,21.4,57.7,21.4,57.7,21.4C57.7,21.4,57.7,21.4,57.7,21.4z M75.3,31.6
	C75.3,31.6,75.3,31.6,75.3,31.6C75.3,31.6,75.3,31.6,75.3,31.6z M75.7,32.3C75.7,32.3,75.7,32.3,75.7,32.3
	C75.7,32.3,75.7,32.3,75.7,32.3z M81.9,25.6c-1.2-1.3-2.8-2.3-6-4.1l-9.7-5.6C63,14,61.3,13,59.6,12.7c-1.5-0.3-3.1-0.3-4.6,0
	c-1.7,0.4-3.3,1.3-6.6,3.2l-9.7,5.6c-3.2,1.9-4.9,2.8-6,4.1c-1,1.2-1.8,2.5-2.3,4c-0.5,1.7-0.5,3.6-0.5,7.3v11.2
	c0,3.8,0,5.6,0.5,7.3c0.5,1.5,1.3,2.9,2.3,4c1.2,1.3,2.8,2.3,6,4.1l9.7,5.6c3.2,1.9,4.9,2.8,6.6,3.2c1.5,0.3,3.1,0.3,4.6,0
	c1.7-0.4,3.3-1.3,6.6-3.2l9.7-5.6c3.2-1.9,4.9-2.8,6-4.1c1-1.2,1.8-2.5,2.3-4c0.5-1.7,0.5-3.6,0.5-7.3V36.9c0-3.8,0-5.6-0.5-7.3
	C83.7,28.1,82.9,26.7,81.9,25.6z"
        />
        <path
          css={css`
            fill-rule: evenodd;
            clip-rule: evenodd;
            fill: #2f80ed;
            animation: ${hideShow} 1.5s ease infinite;
          `}
          d="M46.3,48.1V36.9c0-2,0-3.1-0.1-3.9c0-0.4-0.1-0.6-0.1-0.7c-0.1-0.3-0.2-0.5-0.4-0.7c-0.1,0-0.2-0.2-0.6-0.4
	c-0.7-0.5-1.6-1-3.3-2l-9.7-5.6c-1.7-1-2.7-1.5-3.4-1.9c-0.4-0.2-0.6-0.3-0.6-0.3c-0.3-0.1-0.6-0.1-0.9,0c-0.1,0-0.3,0.1-0.6,0.3
	c-0.7,0.4-1.7,0.9-3.4,1.9l-9.7,5.6c-1.7,1-2.7,1.5-3.3,2c-0.3,0.2-0.5,0.4-0.6,0.4c-0.2,0.2-0.3,0.5-0.4,0.7c0,0.1,0,0.3-0.1,0.7
	c0,0.8-0.1,1.9-0.1,3.9v11.2c0,2,0,3.1,0.1,3.9c0,0.4,0.1,0.6,0.1,0.7c0.1,0.3,0.2,0.5,0.4,0.7c0.1,0,0.2,0.2,0.6,0.4
	c0.7,0.5,1.6,1,3.3,2l9.7,5.6c1.7,1,2.7,1.5,3.4,1.9c0.4,0.2,0.6,0.3,0.6,0.3c0.3,0.1,0.6,0.1,0.9,0c0.1,0,0.3-0.1,0.6-0.3
	c0.7-0.4,1.7-0.9,3.4-1.9l9.7-5.6c1.7-1,2.7-1.5,3.3-2c0.3-0.2,0.5-0.4,0.6-0.4c0.2-0.2,0.3-0.5,0.4-0.7c0-0.1,0-0.3,0.1-0.7
	C46.3,51.2,46.3,50.1,46.3,48.1z M52.3,25.6c-1.2-1.3-2.8-2.3-6-4.1l-9.7-5.6C33.4,14,31.8,13,30,12.7c-1.5-0.3-3.1-0.3-4.6,0
	c-1.7,0.4-3.3,1.3-6.6,3.2l-9.7,5.6c-3.2,1.9-4.9,2.8-6,4.1c-1,1.2-1.8,2.5-2.3,4c-0.5,1.7-0.5,3.6-0.5,7.3v11.2
	c0,3.8,0,5.6,0.5,7.3c0.5,1.5,1.3,2.9,2.3,4c1.2,1.3,2.8,2.3,6,4.1l9.7,5.6c3.2,1.9,4.9,2.8,6.6,3.2c1.5,0.3,3.1,0.3,4.6,0
	c1.7-0.4,3.3-1.3,6.6-3.2l9.7-5.6c3.2-1.9,4.9-2.8,6-4.1c1-1.2,1.8-2.5,2.3-4c0.5-1.7,0.5-3.6,0.5-7.3V36.9c0-3.8,0-5.6-0.5-7.3
	C54.2,28.1,53.4,26.7,52.3,25.6z"
        />
        <path
          css={css`
            fill: #2f80ed;
            animation: ${hideShow} 1.5s ease infinite;
          `}
          d="M42.5,55.5c0.2,0.1,0.4,0.3,0.7,0.4l8,4.6c-1.1,0.9-2.6,1.7-4.9,3.1l-3.8,2.2l-3.8-2.2
	c-2.3-1.4-3.8-2.2-4.9-3.1l8-4.6C42.1,55.7,42.3,55.6,42.5,55.5z"
        />
        <path
          css={css`
            fill: #2d9cdb;
            animation: ${hideShow} 1.5s ease infinite;
          `}
          d="M51.2,24.5c-1.1-0.9-2.6-1.7-4.9-3.1l-3.8-2.2l-3.8,2.2c-2.3,1.4-3.8,2.2-4.9,3.1l8,4.6
	c0.2,0.1,0.5,0.3,0.7,0.4c0.2-0.1,0.4-0.3,0.7-0.4L51.2,24.5z"
        />
      </svg>
    </Container>
  );
};
