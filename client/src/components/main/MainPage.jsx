import React, { useEffect, useState } from 'react'
import PreferenceInput from './PreferenceInput';
import { BarChart3, Brain, ChartSpline, LineChart, ScanSearch } from 'lucide-react';
import ChatDrawer from './ChatDrawer';
import UploadDocument from './UploadDocument';
import api from '../../api/axiosClient';
import RecommendationsSection from '../RecommendationsSection';

export default function MainPage() {

    const [uploadedDocuments, setUploadedDocuments] = useState([]);
    const [recommendations, setRecommendations] = useState([]);
    const [showRecommendation, setShowRecommendation] = useState(false);

    //preference
    const [preference, setPreference] = useState("");

    //upload
    const [uploadedFiles, setUploadedFiles] = useState([])
    const [loading, setLoading] = useState(false)
    const [fileForAnalyze, setFileForAnalyze] = useState();

    const handlePreferenceChange = (e) => {
        setPreference(e.target.value)
    }

    const handleDocumentUpload = (result) => {
        setUploadedDocuments((prev) => [...prev, result.transcriptId]);
        // transform server payload into your UI’s structure
        const picks = (result.recommendation.courses || []).map((c) => ({
            type: "course",
            title: `Course ID ${c.course_id}`,
            course_id: c.course_id,
            description: c.rationale,
            match: c.match,
        }));
        setRecommendations(picks);
        setShowRecommendation(true);
    };

    const handleAnalysis = async () => {
        setLoading(true)
        try {
            const form = new FormData()
            form.append("file", fileForAnalyze)
            const up = await api.post("/transcripts/upload", form, { headers: { "Content-Type": "multipart/form-data" } })
            const transcriptId = up.data.id

            const userPreference = preference;
            // create recommendation
            const reco = await api.post("/recommendations", { transcript_id: transcriptId })
            localStorage.setItem("last_reco_id", reco.data.id);
            // bubble to parent
            handleDocumentUpload({
                transcriptId,
                recommendation: reco.data
            })
        } catch (err) {
            console.error(err)
            alert(err.response?.data?.error || "Upload/analysis failed")
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        const saved = localStorage.getItem("uploadedDocs");
        if (saved) setUploadedDocuments(JSON.parse(saved));
    }, []);


    useEffect(() => {
        localStorage.setItem("uploadedDocs", JSON.stringify(uploadedDocuments));
    }, [uploadedDocuments]);

    return (
        <div className="mx-auto px-4 py-8 sm:px-6 lg:px-8">
            <div className="bg-white border border-gray-200 rounded-2xl shadow-sm p-6">
                <UploadDocument
                    onUpload={handleDocumentUpload}
                    uploadedFiles={uploadedFiles}
                    setUploadedFiles={setUploadedFiles}
                    loading={loading}
                    setLoading={setLoading}
                    setFileForAnalyze={setFileForAnalyze}
                />
                <PreferenceInput value={preference} onChange={handlePreferenceChange} />
                <div className="flex justify-end pb-4">
                    <button
                        className="px-4 py-3 bg-emerald-600 text-white font-medium rounded-xl shadow-sm hover:bg-emerald-700 
                                    hover:shadow-md active:scale-95 transition-all duration-200 inline-flex items-center gap-2
                                max-w-[280px]"
                        onClick={handleAnalysis}
                    >
                        <BarChart3 size={20} strokeWidth={3} />
                        Start Analyzing
                    </button>
                </div>

                {showRecommendation && uploadedDocuments && <RecommendationsSection uploadedDocuments={uploadedDocuments} />}

            </div>

            <ChatDrawer />

        </div>
    )
}
