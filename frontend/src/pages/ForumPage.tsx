import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

const ForumPage: React.FC = () => {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">论坛</h1>
        </div>
        
        <Card>
          <CardHeader>
            <CardTitle>讨论区</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">
              论坛功能正在开发中...
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default ForumPage;