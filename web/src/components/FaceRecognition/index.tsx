

import * as faceapi from "face-api.js";
import React, { useState, useEffect, useRef } from "react";
import { Button, Modal, Progress, Spin, message,Tooltip } from "antd";
import { FormattedMessage, useIntl } from '@umijs/max';
import { QuestionCircleOutlined } from "@ant-design/icons";

interface FaceRecognitionProps {
  visible: boolean;
  onSuccess(array: number[]): void;
  onCancel(): void;
}
const FaceRecognitionModal = (props:FaceRecognitionProps) => {
  const { visible, onSuccess, onCancel } = props;
  const [modelsLoaded, setModelsLoaded] = useState(false);
  const [isCameraCaptured, setIsCameraCaptured] = useState(false);
  const intl = useIntl();
  const videoRef = useRef();
  const canvasRef = useRef();
  const detection = useRef();
  const mediaStreamRef = useRef();
  const [percent, setPercent] = useState(0);
  const loadModels = async () => {
    const MODEL_URL = "/face-recognition";
    Promise.all([
      faceapi.nets.tinyFaceDetector.loadFromUri(MODEL_URL),
      faceapi.nets.faceLandmark68Net.loadFromUri(MODEL_URL),
      faceapi.nets.faceRecognitionNet.loadFromUri(MODEL_URL),
    ]).then((val) => {
      setModelsLoaded(true);
    }).catch((err) => {
      console.log(err);
      message.error(intl.formatMessage({ id: 'pages.personal.settings.faceRecognition.load.models.failed' }));
    });
  };
  useEffect(() => {
    loadModels();
  }, []);
  // 1、chrome://flags/#unsafely-treat-insecure-origin-as-secure
  // 2、选 enabled
  // 3、填写需要调试的 URL，多个 URL 以 , 隔开
  // 4、完全重启 chrome 后起效（改了之后下面也会有个 relaunch 按钮）
  // 5、重启后，本地非https就可以调试
  
  React.useEffect(() => {
    if (visible) {
      setPercent(0);
      if (modelsLoaded) {
        navigator.mediaDevices
          .getUserMedia({ video: { facingMode: "user" } })
          .then((stream) => {
            mediaStreamRef.current = stream;
            setIsCameraCaptured(true);
          }).catch((error) => {
            handleCameraError(error);
          });
      }
    } else {
      clearInterval(detection.current);
      detection.current = null;
      setIsCameraCaptured(false);
    }
    return () => {
      clearInterval(detection.current);
      detection.current = null;
      setIsCameraCaptured(false);
    };
  }, [visible, modelsLoaded]);

  React.useEffect(() => {
    if (isCameraCaptured) {
      let count = 0;
      const interval = setInterval(() => {
        count++;
        if (videoRef.current) {
          videoRef.current.srcObject = mediaStreamRef.current;
          videoRef.current.play();
          clearInterval(interval);
        }
        if (count >= 30) {
          clearInterval(interval);
          onCancel();
        }
      }, 100);
    } else {
      mediaStreamRef.current?.getTracks().forEach(track => track.stop());
      if (videoRef.current) {
        videoRef.current.srcObject = null;
      }
    }
  }, [isCameraCaptured]);

  const handleStreamVideo = () => {
    let count = 0;
    let goodCount = 0;
    if (!detection.current) {
      detection.current = setInterval(async () => {
        if (modelsLoaded && videoRef.current && visible) {
          const faces = await faceapi.detectAllFaces(videoRef.current, new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptors();

          count++;
          if (count % 50 === 0) {
            message.warning(intl.formatMessage({ id: 'pages.personal.settings.faceRecognition.face.data.extract.tips' }));

          } else if (count > 300) {
            message.warning(intl.formatMessage({ id: 'pages.personal.settings.faceRecognition.face.data.extract.failed' }));
            onCancel();
          }
          if (faces.length === 1) {
            const face = faces[0];
            setPercent(Math.round(face.detection.score * 100));
            const array = Array.from(face.descriptor);
            if (face.detection.score > 0.9) {
              goodCount++;
              if (face.detection.score > 0.99 || goodCount > 10) {
                clearInterval(detection.current);
                onSuccess(array);
              }
            }
          } else {
            setPercent(Math.round(percent / 2));
          }
        }
      }, 100);
    }
  };
  const handleCameraError = (error:DOMException) => {
    onCancel();
    if (error instanceof DOMException) {
      if (error.name === "NotFoundError" || error.name === "DevicesNotFoundError") {
      message.error(intl.formatMessage({ id: 'pages.login.face.recognition.device.not.found' }));
      } else if (error.name === "NotAllowedError" || error.name === "PermissionDeniedError") {
      message.error(intl.formatMessage({ id: 'pages.login.face.recognition.device.not.permission' }));
      } else if (error.name === "NotReadableError" || error.name === "TrackStartError") {
      message.error(intl.formatMessage({ id: 'pages.login.face.recognition.device.used.by.other.app' }));
      } else if (error.name === "TypeError") {
      message.error(intl.formatMessage({ id: 'pages.login.face.recognition.device.protocol.need.https' }));
      } else {
      message.error(intl.formatMessage({ id: 'pages.unknown.error' }));
      }
    }
  };

  return (
    <div>
      <Modal
        closable={false}
        maskClosable={false}
        destroyOnClose={true}
        open={visible && isCameraCaptured}
        title={<>
        <FormattedMessage id='pages.faceRecognition.face.data.extract'  />
        <Tooltip title={<FormattedMessage id="pages.faceRecognition.operation.description"  />}>
              <QuestionCircleOutlined />
            </Tooltip>  
        </>}
        width={350}
        footer={[
          <Button key="back" onClick={onCancel}>
           <FormattedMessage id='pages.operation.cancel' />
          </Button>,
        ]}
      >
        <Progress percent={percent} />
        <div style={{ marginTop: "20px", marginBottom: "50px", justifyContent: "center", alignContent: "center", position: "relative", flexDirection: "column" }}>
          {
            modelsLoaded ?
              <div style={{ display: "flex", justifyContent: "center", alignContent: "center" }}>
                <video
                  ref={videoRef}
                  onPlay={handleStreamVideo}
                  style={{
                    borderRadius: "50%",
                    height: "220px",
                    verticalAlign: "middle",
                    width: "220px",
                    objectFit: "cover",
                  }}
                ></video>
                <div style={{
                  position: "absolute",
                  width: "240px",
                  height: "240px",
                  top: "50%",
                  left: "50%",
                  transform: "translate(-50%, -50%)",
                }}>
                  <svg width="240" height="240" fill="none">
                    <circle
                      strokeDasharray="700"
                      strokeDashoffset={700 - 6.9115 * percent}
                      strokeWidth="4"
                      cx="120"
                      cy="120"
                      r="110"
                      stroke="#32CD32"
                      transform="rotate(-90, 120, 120)"
                      strokeLinecap="round"
                      style={{ transition: "all .2s linear" }}
                    ></circle>
                  </svg>
                </div>
                <canvas ref={canvasRef} style={{ position: "absolute" }} />
              </div>
              :
              <div>
                <Spin tip={intl.formatMessage({id:'pages.loading'})} size="large" style={{ display: "flex", justifyContent: "center", alignContent: "center" }}>
                  <div className="content" />
                </Spin>
              </div>
          }
        </div>
      </Modal>
    </div>
  );
};

export default FaceRecognitionModal;
