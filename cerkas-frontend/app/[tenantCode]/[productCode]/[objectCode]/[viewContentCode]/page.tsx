"use client";

import { Card, CardContent } from "@/components/ui/card";
import { toLabel } from "@/lib/utils";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

interface RouteParams {
  tenantCode: string;
  productCode: string;
  objectCode: string;
  viewContentCode: string;
  [key: string]: string | string[];
}

interface Field {
  field_code: string;
  field_name: string;
}

const DynamicTable = ({ fields }: { fields: Field[] }) => {
  return (
    <div className="overflow-x-auto">
      <table className="table-auto border border-gray-100 whitespace-nowrap">
        <thead>
          <tr className="bg-gray-100">
            {fields.map((field) => (
              <th
                key={field.field_code}
                className="px-2 py-2 border border-gray-100 text-left"
                style={{
                  minWidth: `${field.field_name.length * 10 + 40}px`, // dynamic width
                }}
              >
                {field.field_name}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          <tr>
            {fields.map((field) => (
              <td
                key={field.field_code}
                className="px-2 py-2 border border-gray-100 text-gray-600"
                style={{
                  minWidth: `${field.field_name.length * 10 + 40}px`, // match header width
                }}
              >
                {/* Placeholder */} example
              </td>
            ))}
          </tr>
        </tbody>
      </table>
    </div>
  );
};


export default function DynamicPage() {
  const params = useParams<RouteParams>();
  const { tenantCode, productCode, objectCode, viewContentCode } = params;

  const [responseData, setResponseData] = useState<any>(null); // You can define a more specific type if you know the shape
  const [viewContent, setViewContent] = useState<any>(null);
  const [viewLayout, setViewLayout] = useState<any>(null);

  const [isAPIResponseAccordionOpen, setIsAPIResponseAccordionOpen] = useState(false);
  const [isDynamicParamAccordionOpen, setIsDynamicParamAccordionOpen] = useState(false);

  useEffect(() => {
    const sendRequest = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/t/${tenantCode}/p/${productCode}/o/${objectCode}/view/${viewContentCode}/record`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({}), // Add any payload here if needed
          }
        );

        if (!response.ok) {
          throw new Error("Failed to fetch");
        }

        const data = await response.json();
        setResponseData(data.data);
        setViewContent(data.data.view_content);
        setViewLayout(data.data.layout);
      } catch (error) {
        console.error("API error:", error);
        setResponseData({ error: (error as Error).message });
      }
    };

    if (tenantCode && productCode && objectCode && viewContentCode) {
      sendRequest();
    }
  }, [tenantCode, productCode, objectCode, viewContentCode]);

  type ViewChild = {
    type: string;
    class_name?: string;
    props?: {
      fields?: any[]; // you can define a stricter type if needed
    };
  };


  return (
    <div className="flex flex-col items-left justify-left min-h-screen bg-gray-100">
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold text-cyan-600 mb-3">
          {viewContent?.object?.display_name ? viewContent?.object?.display_name : toLabel(objectCode)}
        </h1>

        <div className="space-y-4">
          {viewLayout?.children.map((child: ViewChild, index: number) => {
            if (child.type === "table") {
              return (
                <Card key={index} className="shadow-md">
                  <CardContent className="p-2 pb-8 overflow-x-auto">
                    <DynamicTable fields={child.props?.fields || []} />
                  </CardContent>
                </Card>
              );
            }
            return null;
          })}
        </div>

        <div className="mt-4">
          <button
            onClick={() => setIsDynamicParamAccordionOpen(!isDynamicParamAccordionOpen)}
            className="w-full flex justify-between items-center text-left bg-indigo-100 px-4 py-2 rounded-lg font-semibold text-gray-800 hover:bg-indigo-200 transition"
          >
            <span>Dynamic Param</span>
            <span>{isDynamicParamAccordionOpen ? "−" : "+"}</span>
          </button>

          {isDynamicParamAccordionOpen && (
            <pre className="mt-2 bg-gray-200 text-sm p-4 rounded overflow-auto text-gray-800">
              <div className="space-y-3 text-gray-700">
                <p>
                  <span className="font-semibold">Tenant Code:</span>{" "}
                  <span className="text-indigo-600">{tenantCode}</span>
                </p>
                <p>
                  <span className="font-semibold">Product Code:</span>{" "}
                  <span className="text-indigo-600">{productCode}</span>
                </p>
                <p>
                  <span className="font-semibold">Object Code:</span>{" "}
                  <span className="text-indigo-600">{objectCode}</span>
                </p>
                <p>
                  <span className="font-semibold">View Content Code:</span>{" "}
                  <span className="text-indigo-600">{viewContentCode}</span>
                </p>
              </div>
            </pre>
          )}
        </div>

        <div className="mt-4">
          <button
            onClick={() => setIsAPIResponseAccordionOpen(!isAPIResponseAccordionOpen)}
            className="w-full flex justify-between items-center text-left bg-indigo-100 px-4 py-2 rounded-lg font-semibold text-gray-800 hover:bg-indigo-200 transition"
          >
            <span>Layout API Response</span>
            <span>{isAPIResponseAccordionOpen ? "−" : "+"}</span>
          </button>

          {isAPIResponseAccordionOpen && (
            <SyntaxHighlighter
              language="json"
              style={vscDarkPlus}
              showLineNumbers
              wrapLines
              customStyle={{
                borderRadius: '0.5rem',
                padding: '1rem',
                fontSize: '0.875rem',
                backgroundColor: '#1e1e1e'
              }}
            >
              {JSON.stringify(responseData, null, 2)}
            </SyntaxHighlighter>
          )}
        </div>
      </div>
    </div>
  );
}
