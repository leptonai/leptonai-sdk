import { FC, PropsWithChildren, ReactNode } from "react";
import { useAntdTheme } from "@lepton-dashboard/hooks/use-antd-theme";
import { css } from "@emotion/react";
import { EmotionProps } from "@lepton-dashboard/interfaces/emotion-props";
import { Skeleton } from "antd";

export const Card: FC<
  PropsWithChildren<
    {
      title?: ReactNode;
      titleIcon?: ReactNode;
      loading?: boolean;
      extra?: ReactNode;
      borderless?: boolean;
      shadowless?: boolean;
      paddingless?: boolean;
    } & EmotionProps
  >
> = ({
  children,
  title,
  titleIcon,
  extra,
  loading = false,
  borderless = false,
  shadowless = false,
  paddingless = false,
  className,
}) => {
  const theme = useAntdTheme();
  return (
    <div
      css={css`
        display: flex;
        flex-direction: column;
        overflow: hidden;
        background-color: ${theme.colorBgContainer};
        border-color: ${theme.colorBorder};
        border-radius: ${theme.borderRadius}px;
        border-style: solid;
        border-width: ${borderless ? 0 : "1px"};
        box-shadow: ${shadowless ? "none" : theme.boxShadowTertiary};
      `}
    >
      {(title || extra) && (
        <div
          css={css`
            flex: 0 0 48px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            border-bottom: 1px solid ${theme.colorBorder};
            background-color: ${theme.colorBgContainer};
            color: ${theme.colorTextHeading};
            font-size: 14px;
            font-weight: 500;
            padding: 0 16px;
          `}
        >
          <div
            css={css`
              display: flex;
            `}
          >
            {titleIcon && (
              <div
                css={css`
                  margin-right: 12px;
                `}
              >
                {titleIcon}
              </div>
            )}
            {title}
          </div>
          <div>{extra}</div>
        </div>
      )}
      <div
        className={className}
        css={css`
          padding: ${paddingless ? 0 : "16px"};
        `}
      >
        {loading ? <Skeleton /> : children}
      </div>
    </div>
  );
};
