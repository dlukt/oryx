//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";
import {useNavigate} from "react-router-dom";
import {Errors, Token} from "../utils";
import axios from "axios";

export default function useDvrVodStatus() {
  const navigate = useNavigate();
  const [status, setStatus] = React.useState({ vod: undefined, dvr: undefined });

  React.useEffect(() => {
    const p0 = axios.post('/terraform/v1/hooks/vod/query', {}, {
      headers: Token.loadBearerHeader(),
    });
    const p1 = axios.post('/terraform/v1/hooks/dvr/query', {}, {
      headers: Token.loadBearerHeader(),
    });

    Promise.all([p0, p1]).then(([vodRes, dvrRes]) => {
      console.log(`VodPattern: Query ok, ${JSON.stringify(vodRes.data.data)}`);
      console.log(`DvrPattern: Query ok, ${JSON.stringify(dvrRes.data.data)}`);
      // Update state once to prevent unnecessary re-renders.
      setStatus({ vod: vodRes.data.data, dvr: dvrRes.data.data });
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate]);

  return [status.dvr, status.vod];
}

export function useRecordStatus() {
  const navigate = useNavigate();
  const [recordStatus, setRecordStatus] = React.useState();

  React.useEffect(() => {
    axios.post('/terraform/v1/hooks/record/query', {
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      console.log(`RecordPattern: Query ok, ${JSON.stringify(res.data.data)}`);
      setRecordStatus(res.data.data);
    }).catch(e => {
      const err = e.response.data;
      if (err.code === Errors.auth) {
        alert(`Token过期，请重新登录，${err.code}: ${err.data.message}`);
        navigate('/routers-logout');
      } else {
        alert(`服务器错误，${err.code}: ${err.data.message}`);
      }
    });
  }, [navigate]);

  return recordStatus;
}

