import React from "react";
import IconButton from "./IconButton";
import * as Icon from 'react-bootstrap-icons';
import {Clipboard} from "../utils";
import {useTranslation} from "react-i18next";

export default function CopyButton({text, title, className, size = 20}) {
  const [copied, setCopied] = React.useState(false);
  const {t} = useTranslation();

  const handleCopy = React.useCallback((e) => {
    e.preventDefault();
    Clipboard.copy(text).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }).catch((err) => {
      alert(`${t('helper.copyFail')} ${err}`);
    });
  }, [text, t]);

  const label = copied ? t('helper.copyOk') : (title || t('helper.copy'));

  return (
    <IconButton title={label} onClick={handleCopy} className={className}>
      {copied ? <Icon.Check size={size} className="text-success"/> : <Icon.Clipboard size={size}/>}
    </IconButton>
  );
}
