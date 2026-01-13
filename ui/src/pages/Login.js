//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";
import Container from "react-bootstrap/Container";
import {Form, Button, Spinner, InputGroup} from 'react-bootstrap';
import axios from "axios";
import {useNavigate} from "react-router-dom";
import {Token, Tools} from '../utils';
import {SrsErrorBoundary} from "../components/SrsErrorBoundary";
import {useErrorBoundary} from "react-error-boundary";
import {useTranslation} from "react-i18next";
import {Eye, EyeSlash} from "react-bootstrap-icons";

export default function Login({onLogin}) {
  return (
    <SrsErrorBoundary>
      <LoginImpl onLogin={onLogin} />
    </SrsErrorBoundary>
  );
}

function LoginImpl({onLogin}) {
  const [plaintext, setPlaintext] = React.useState(false);
  const [password, setPassword] = React.useState('');
  const [operating, setOperating] = React.useState(false);
  const navigate = useNavigate();
  // We use a single ref for the password input now.
  const passwordRef = React.useRef();
  const { showBoundary: handleError } = useErrorBoundary();
  const {t} = useTranslation();

  // Verify the token if exists.
  React.useEffect(() => {
    const token = Token.load();
    if (!token || !token.token) return;

    console.log(`Login: Verify, token is ${Tools.mask(token)}`);

    // Both JWT token or Bearer token are OK. Here we use JWT token.
    axios.post('/terraform/v1/mgmt/token', {
      ...token,
    }).then(res => {
      // Here we use the Bearer token to verify again.
      axios.post('/terraform/v1/mgmt/token', {
      }, {
        headers: Token.loadBearerHeader(),
      }).then(res => {
        console.log(`Login: Done, token is ${Tools.mask(token)}`);
        navigate('/routers-scenario');
      });
    }).catch(handleError);
  }, [navigate, handleError]);

  // Focus to password input on mount or when visibility changes (optional, but good for UX)
  // Actually with single component, focus should remain if we just toggle type,
  // but if we want to ensure focus, we can keep this.
  // However, toggling type shouldn't lose focus in modern browsers.
  // Let's remove the manual focus effect that was needed for component swapping.
  // We'll focus only on mount if needed, or rely on autoFocus if we added it (we didn't).

  // User click login button.
  const handleLogin = React.useCallback((e) => {
    e.preventDefault();
    setOperating(true);

    axios.post('/terraform/v1/mgmt/login', {
      password,
    }).then(async (res) => {
      await new Promise(resolve => setTimeout(resolve, 600));

      const data = res.data.data;
      console.log(`Login: OK, token is ${Tools.mask(data)}`);
      Token.save(data);

      onLogin && onLogin();
      navigate('/routers-scenario');
    }).catch(handleError).finally(setOperating);
  }, [password, handleError, onLogin, navigate, setOperating]);

  return (
    <>
      <Container fluid>
        <Form>
          <Form.Group className="mb-3" controlId="formBasicPassword">
            <Form.Label>{t('login.passwordLabel')}</Form.Label>
            <InputGroup>
              <Form.Control
                type={plaintext ? "text" : "password"}
                placeholder="Password"
                ref={passwordRef}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <Button
                variant="outline-secondary"
                onClick={() => setPlaintext(!plaintext)}
                onMouseDown={(e) => e.preventDefault()}
                aria-label={plaintext ? "Hide password" : "Show password"}
                title={plaintext ? "Hide password" : "Show password"}
              >
                {plaintext ? <EyeSlash /> : <Eye />}
              </Button>
            </InputGroup>
            <Form.Text className="text-muted">
              * {t('login.passwordTip')}
            </Form.Text>
          </Form.Group>
          {/* Checkbox removed in favor of InputGroup button */}
          <Button variant="primary" type="submit" disabled={operating} onClick={(e) => handleLogin(e)}>
            {operating && (
              <Spinner
                as="span"
                animation="border"
                size="sm"
                role="status"
                aria-hidden="true"
                className="me-2"
              />
            )}
            {t('login.labelLogin')}
          </Button>
        </Form>
      </Container>
    </>
  );
}
